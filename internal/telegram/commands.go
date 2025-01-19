package telegram

import "log"

type Command func(WebhookMessage)

func CommandStart(message WebhookMessage) {
	text := `
Olá! Eu sou o assistente virtual independente, da Cidade de São José dos Campos.

Comandos disponíveis:
/start, /ajuda - Inicia a conversa com o assistente
/sobre - Exibe informações sobre o bot

# CETESB
/cetesb:qualidade-ar - Exibe a qualidade do ar na cidade
	`

	err := SendMessage(message.Chat.ID, text)
	if err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

func CommandAbout(message WebhookMessage) {
	text := `Este bot foi desenvolvido por https://GustavoNovaes.dev, para ajudar a população de São José dos Campos a obter informações sobre a cidade de forma facilitada e automatizada a partir dos serviços de chat Telegram e Discord.`
	err := SendMessage(message.Chat.ID, text)
	if err != nil {
		log.Printf("failed to send message: %v", err)
	}
}
