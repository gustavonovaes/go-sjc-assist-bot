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

func HandleWebhook(w http.ResponseWriter, r *http.Request, commands map[string]Command) error {
	w.WriteHeader(http.StatusOK)

	var webhookResponse WebhookResponse

	if err := json.NewDecoder(r.Body).Decode(&webhookResponse); err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}
	defer r.Body.Close()

	log.Printf("Received message: %+v", webhookResponse.Message)

	if len(commands) == 0 {
		log.Println("No commands available")
	}

	for command, handler := range commands {
		if strings.Contains(webhookResponse.Message.Text, command) {
			log.Printf(
				"User %s requested command %s",
				webhookResponse.Message.From.Username,
				command,
			)

			handler(webhookResponse.Message)
		}
	}

	return nil
}

func SendMessage(chatID int, message string) error {
	res, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_TOKEN+"/sendMessage",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"chat_id": %d, "text": "%s"}`, chatID, message)),
	)

	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", res.StatusCode)
	}

	return nil
}
