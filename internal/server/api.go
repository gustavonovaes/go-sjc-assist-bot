package server

import (
	"net/http"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
)

func startTelegramWebhookServer() {
	telegram.SetupWebhook()

	// Start the server
	http.HandleFunc("/telegram", telegram.HandleWebhook)

}
