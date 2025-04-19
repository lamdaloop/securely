package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lamdaloop/securely/handlers"
	"github.com/lamdaloop/securely/utils"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found. Falling back to system env.")
	}

	// Ensure the secrets directory exists
	if _, err := os.Stat("./secrets"); os.IsNotExist(err) {
		if err := os.Mkdir("./secrets", 0700); err != nil {
			log.Fatalf("Failed to create secrets directory: %v", err)
		}
	}

	// Start the background cleaner
	utils.StartCleaner(1 * 60 * 1e9) // 1 minute

	// Route handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/secret/") {
			http.ServeFile(w, r, "static/view.html")
		} else {
			http.ServeFile(w, r, "static/index.html")
		}
	})

	http.HandleFunc("/auth/login", handlers.LoginHandler)
	http.HandleFunc("/auth/callback", handlers.CallbackHandler)
	http.HandleFunc("/auth/logout", handlers.LogoutHandler)
	http.HandleFunc("/auth/me", handlers.WhoAmI)
	http.HandleFunc("/api/secret", handlers.RequireAuth(handlers.HandleSecret))
	http.HandleFunc("/api/secret/", handlers.RequireAuth(handlers.HandleRetrieveSecret))
	http.HandleFunc("/api/retrieve/", handlers.HandleRetrieveSecret)


	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🚀 Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
