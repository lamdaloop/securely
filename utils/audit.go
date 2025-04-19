package utils

import (
	"fmt"
	"os"
	"time"
)

func WriteAudit(event, user, secretID string) {
	logLine := fmt.Sprintf("%s | event=%s | user=%s | secret=%s\n",
		time.Now().Format(time.RFC3339), event, user, secretID)

	f, err := os.OpenFile("logs/audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ audit log write failed:", err)
		return
	}
	defer f.Close()
	f.WriteString(logLine)
}
