package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/lamdaloop/securely/models"
	"github.com/lamdaloop/securely/storage"
	"github.com/lamdaloop/securely/utils"
)

type createSecretRequest struct {
	Message         string `json:"message"`
	ExpireInMinutes int    `json:"expire_in_minutes"`
	OneTime         bool   `json:"one_time"`
	Password        string `json:"password"`
}

type createSecretResponse struct {
	ID string `json:"id"`
}

type retrieveSecretRequest struct {
	Password string `json:"password"`
}

type retrieveSecretResponse struct {
	Message   string `json:"message"`
	CreatedBy string `json:"created_by"`
}

func HandleSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req createSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := models.GenerateID()
	encrypted, iv, err := models.Encrypt([]byte(req.Message))
	if err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	user := getUserEmail(r)
	secret := models.Secret{
		ID:           id,
		EncryptedMsg: encrypted,
		IV:           iv,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(time.Duration(req.ExpireInMinutes) * time.Minute),
		OneTime:      req.OneTime,
		CreatedBy:    user,
	}

	if req.Password != "" {
		hash, err := models.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Password hashing failed", http.StatusInternalServerError)
			return
		}
		secret.PasswordHash = hash
	}

	if err := storage.SaveSecret(secret); err != nil {
		http.Error(w, "Failed to save secret", http.StatusInternalServerError)
		return
	}

	utils.WriteAudit("created", user, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createSecretResponse{ID: id})
}

func HandleRetrieveSecret(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	var req retrieveSecretRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	secret, err := storage.LoadSecret(id)
	if err != nil {
		log.Println("❌ Load failed:", err)
		http.Error(w, "Secret not found", http.StatusNotFound)
		return
	}

	if time.Now().After(secret.ExpiresAt) {
		_ = storage.DeleteSecret(id)
		utils.WriteAudit("expired", getUserEmail(r), id)
		http.Error(w, "Secret expired", http.StatusGone)
		return
	}

	if secret.PasswordHash != "" {
		if req.Password == "" || !models.CheckPasswordHash(req.Password, secret.PasswordHash) {
			utils.WriteAudit("unauthorized", getUserEmail(r), id)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	plaintext, err := models.Decrypt(secret.EncryptedMsg, secret.IV)
	if err != nil {
		utils.WriteAudit("decryption_failed", getUserEmail(r), id)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	if secret.OneTime {
		_ = storage.DeleteSecret(id)
	}

	utils.WriteAudit("accessed", getUserEmail(r), id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retrieveSecretResponse{
		Message:   string(plaintext),
		CreatedBy: secret.CreatedBy,
	})
}

func getUserEmail(r *http.Request) string {
	cookie, err := r.Cookie("user_email")
	if err != nil {
		return "unknown"
	}
	return cookie.Value
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/secret/") {
		http.ServeFile(w, r, "static/view.html")
	} else {
		http.ServeFile(w, r, "static/index.html")
	}
}
