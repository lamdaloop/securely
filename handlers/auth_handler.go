package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleapi "google.golang.org/api/oauth2/v2"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	stateToken = "secretbox-oauth-state"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(stateToken)
	http.Redirect(w, r, url, http.StatusFound)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != stateToken {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	service, err := googleapi.New(client)
	if err != nil {
		http.Error(w, "Failed to create Google API client", http.StatusInternalServerError)
		return
	}

	userinfo, err := service.Userinfo.Get().Do()
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	if !strings.HasSuffix(userinfo.Email, "@gmail.com") {
		http.Error(w, "Access denied: outside organization", http.StatusForbidden)
		return
	}

	// Set session cookie with 1-hour expiration
	http.SetCookie(w, &http.Cookie{
		Name:     "user_email",
		Value:    userinfo.Email,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600, // 1 hour
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user_email",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func WhoAmI(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_email")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not logged in"})
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{"email": cookie.Value})
}
