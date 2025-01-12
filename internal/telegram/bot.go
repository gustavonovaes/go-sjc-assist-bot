// Package telegram fornece funcionalidades para acessar a API do Telegram.
package telegram

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"

	"gustavonovaes.dev/go-sjc-assist-bot/pkg/config"
)

var (
	appConfig          config.Config
	webhookSecretToken string
)

func init() {
	webhookSecretToken = randomString(32)
	appConfig = config.New()
}

func SetupWebhook() error {
	res, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_TOKEN+"/setWebhook?url="+appConfig.TELEGRAM_WEBHOOK_URL+"&secret_token="+webhookSecretToken,
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

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling webhook")
	w.WriteHeader(http.StatusOK)

	requestSecretToken := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if requestSecretToken != webhookSecretToken {
		log.Println("Invalid secret token")
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body")
		return
	}

	log.Println("Request body: ", string(content))
}

func randomString(length int) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
