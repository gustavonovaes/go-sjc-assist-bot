// Package telegram fornece funcionalidades para acessar a API do Telegram.
package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gustavonovaes.dev/go-sjc-assist-bot/pkg/config"
)

var (
	appConfig config.Config
)

func init() {
	appConfig = config.New()
}

func SetupWebhook() error {
	res, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_TOKEN+"/setWebhook?url="+appConfig.TELEGRAM_WEBHOOK_URL,
		"application/json",
		nil,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyContent := make([]byte, 256)
	res.Body.Read(bodyContent)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to setup webhook, status code: %d, body: %s",
			res.StatusCode,
			bodyContent,
		)
	}

	return nil
}

func HandleWebhook(w http.ResponseWriter, r *http.Request, commands map[string]Command) {
	log.Println("Handling webhook")
	w.WriteHeader(http.StatusOK)

	var webhookMessage WebhookMessage
	if err := json.NewDecoder(r.Body).Decode(&webhookMessage); err != nil {
		log.Println("Failed to decode request body")
		return
	}

	log.Printf("Received message: %+v", webhookMessage)

	if len(commands) == 0 {
		log.Println("No commands available")
		return
	}

	for command, handler := range commands {
		if strings.Contains(webhookMessage.Text, command) {
			log.Printf("User %s requested command %s", webhookMessage.From.Username, command)
			handler(webhookMessage)
			return
		}
	}
}

func SendMessage(chatID int, message string) {
	res, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_TOKEN+"/sendMessage",
		"application/json",
		nil,
	)

	if err != nil {
		log.Println("Failed to send message")
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("Failed to send message")
		return
	}
}
