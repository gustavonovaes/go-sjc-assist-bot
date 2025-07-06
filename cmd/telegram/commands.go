package main

import (
	"fmt"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/sspsp"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/infra/telegram"
)

const QUALAR_STATION_ID = 49 // QUALAR_STATION_ID is the default station ID for São José dos Campos
const MUNICIPALITY_ID = 560  // MUNICIPALITY_ID is the municipality ID for São José dos Campos

func CommandQualityAir(message *telegram.WebhookMessage) error {
	commandCityId := ExtractCityIdFromMessage(message.Text)
	res, err := cetesb.GetQualarData(commandCityId)
	if err != nil {
		return fmt.Errorf(
			"failed to get data from CETESB: %v\n%v",
			err,
			telegram.SendMessage(message.Chat.ID, "Erro ao obter dados da CETESB"),
		)
	}

	return telegram.SendMessage(
		message.Chat.ID,
		fmt.Sprintf(
			`<b>Indice qualidade do Ar:</b>\n%.0f - %s\n`,
			res.Features[0].Attributes.Indice,
			res.Features[0].Attributes.Qualidade,
		),
	)
}

func CommandCrimes(message *telegram.WebhookMessage) error {
	data, err := sspsp.GetPoliceIncidentsCriminal(MUNICIPALITY_ID)
	if err != nil {
		return fmt.Errorf(
			"failed to get data from SSP-SP: %v\n%v",
			err,
			telegram.SendMessage(message.Chat.ID, "Erro ao obter dados da SSP-SP"),
		)
	}

	text := sspsp.GenerateCrimeStatisticsTable(data[:10])

	return telegram.SendMessage(
		message.Chat.ID,
		fmt.Sprintf(
			"<b>Crimes registrados em São José dos Campos:</b>\n%s",
			text,
		),
	)
}

func CommandStart(message *telegram.WebhookMessage) error {
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

	return telegram.SendMessage(message.Chat.ID, text)
}

func CommandAbout(message *telegram.WebhookMessage) error {
	text := `Este bot foi desenvolvido por https://GustavoNovaes.dev, para ajudar a população de São José dos Campos a obter informações sobre a cidade de forma facilitada e automatizada a partir dos serviços de chat Telegram e Discord.`
	return telegram.SendMessage(message.Chat.ID, text)
}
