package telegram

type Command func(WebhookMessage) error

func CommandStart(message WebhookMessage) error {
	text := `
Olá! Eu sou o assistente virtual independente, da Cidade de São José dos Campos.

Comandos disponíveis:
/start, /ajuda - Inicia a conversa com o assistente
/sobre - Exibe informações sobre o bot

<b>CETESB</b>
/qualidadeAr - Exibe a qualidade do ar na cidade

<b>SSP-SP</b>
/crimes - Exibe os crimes registrados na cidade
	`

	return SendMessage(message.Chat.ID, text)
}

func CommandAbout(message WebhookMessage) error {
	text := `Este bot foi desenvolvido por https://GustavoNovaes.dev, para ajudar a população de São José dos Campos a obter informações sobre a cidade de forma facilitada e automatizada a partir dos serviços de chat Telegram e Discord.`
	return SendMessage(message.Chat.ID, text)
}
