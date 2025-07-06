package main

import (
	"fmt"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/sspsp"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/infra/telegram"
)

const QUALAR_STATION_ID = 49 // QUALAR_STATION_ID is the default station ID for São José dos Campos
const MUNICIPALITY_ID = 560  // MUNICIPALITY_ID is the municipality ID for São José dos Campos

func CommandStart(message *telegram.WebhookMessage) error {
	text := `
👋 <b>Bem-vindo(a)!</b>
Eu sou o assistente virtual da Cidade de São José dos Campos, aqui para fornecer informações úteis e práticas sobre a cidade.

<b>Comandos disponíveis:</b>
/start, /ajuda - Inicia a conversa com o bot
/sobre - Exibe informações sobre o bot

<b>🌱 CETESB</b>
/qualidadeAr - Exibe o índice de qualidade do ar da cidade via CETESB

<b>🚔 SSP-SP</b>
/crimes - Exibe o total de crimes registrados na cidade nos últimos anos
/mapaCrimes - Exibe link para o mapa com as marcações dos crimes registrados no último semestre
	`

	return telegram.SendMessage(message.Chat.ID, text)
}

func CommandAbout(message *telegram.WebhookMessage) error {
	text := `Este bot foi desenvolvido por <a href="https://gustavonovaes.dev">Gustavo Novaes</a> para auxiliar a população de São José dos Campos. Ele fornece informações úteis sobre a cidade de forma prática e automatizada, utilizando os serviços de chat do Telegram e Discord.

Se você tiver sugestões ou encontrar problemas, entre em contato ou contribua no repositório do projeto no GitHub: <a href="https://github.com/GustavoNovaes/go-sjc-assist-bot">github.com/GustavoNovaes/go-sjc-assist-bot</a>.`
	return telegram.SendMessage(message.Chat.ID, text)
}

func CommandQualityAir(message *telegram.WebhookMessage) error {
	commandCityId := ExtractCityIdFromMessage(message.Text)
	if commandCityId == 0 {
		commandCityId = QUALAR_STATION_ID
	}

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

	text := fmt.Sprintf(
		"<b> Total de Crimes nos últimos Anos - São José dos Campos:</b>\n <code>%s</code>",
		sspsp.GenerateCrimeStatisticsTable(data[:10]),
	)

	return telegram.SendMessage(message.Chat.ID, text)
}

func CommandMapCrimes(message *telegram.WebhookMessage) error {
	text := `
<b>🗺️ Mapa de Crimes - São José dos Campos</b>
Mapa com marcações dos crimes registrados na cidade no primeiro semestre de 2025.

Link para o mapa: 
https://www.google.com/maps/d/viewer?mid=1Z-LoxrmX55O5_Odo1lRXoCcs5TOXifs

Os dados criminais foram obtidos através do <b>Portal Transparência - Números sem Mistério</b> da SSP-SP. Link para o portal:
https://www.ssp.sp.gov.br/estatistica/consultas
	`

	return telegram.SendMessage(message.Chat.ID, text)
}
