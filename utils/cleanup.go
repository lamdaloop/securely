package utils

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/lamdaloop/securely/storage"
)

func StartCleaner(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			files, err := os.ReadDir("secrets")
			if err != nil {
				log.Println("Cleanup: failed to read secrets dir:", err)
				continue
			}

			for _, f := range files {
				if !strings.HasSuffix(f.Name(), ".bin") {
					continue
				}
				id := strings.TrimSuffix(f.Name(), ".bin")
				secret, err := storage.LoadSecret(id)
				if err != nil {
					continue
				}
				if time.Now().After(secret.ExpiresAt) {
					_ = storage.DeleteSecret(id)
					log.Printf("🧹 Deleted expired secret: %s\n", id)
				}
			}
		}
	}()
}
