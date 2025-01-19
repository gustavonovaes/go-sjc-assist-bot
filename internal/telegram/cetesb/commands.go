package cetesb

import "gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"

func CommandQualidadeAr(message telegram.WebhookMessage) {
	telegram.SendMessage(message.MessageID, "Comando /cetesb:qualidade-ar")
}
