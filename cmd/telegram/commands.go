package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/news"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/nlp"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/sspsp"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/infra/telegram"
)

const QUALAR_STATION_ID = 49 // QUALAR_STATION_ID is the default station ID for São José dos Campos
const MUNICIPALITY_ID = 560  // MUNICIPALITY_ID is the municipality ID for São José dos Campos

func CommandStart(message *telegram.WebhookMessage) error {
	text := `👋 Olá! Eu sou um assistente virtual da Cidade de São José dos Campos. Estou aqui para te ajudar com algumas informações sobre a cidade.

*Comandos disponíveis:*
- /start, /ajuda: Inicia a conversa com o bot.
- /sobre: Exibe informações sobre o bot.

*🌱 CETESB*
- /qualidadeAr: Exibe o índice de qualidade do ar da cidade via CETESB.

*🚔 SSP-SP*
- /crimes: Exibe o total de crimes registrados na cidade nos últimos anos.
- /mapaCrimes: Exibe link para o mapa com as marcações dos crimes registrados no último semestre.

*📰 Notícias*
- /ultimasNoticias: Exibe as últimas notícias da cidade dos principais portais.
	`

	return telegram.SendMessage(message.Chat.ID, text)
}

func CommandAbout(message *telegram.WebhookMessage) error {
	text := `Este bot foi desenvolvido por [Gustavo Novaes](https://gustavonovaes.dev) para auxiliar a população de São José dos Campos. Ele fornece informações úteis sobre a cidade de forma prática e automatizada, utilizando os serviços de chat do Telegram e Discord.

Se você tiver sugestões ou encontrar problemas, entre em contato ou contribua no repositório do projeto no GitHub: [github.com/GustavoNovaes/go-sjc-assist-bot](https://github.com/GustavoNovaes/go-sjc-assist-bot).`
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
			`*Indice qualidade do Ar:*\n%.0f - %s\n`,
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
		"*Total de Crimes nos últimos anos - São José dos Campos:*\n ```\n%s\n```",
		sspsp.GenerateCrimeStatisticsTable(data[:10]),
	)

	return telegram.SendMessage(message.Chat.ID, text)
}

func CommandMapCrimes(message *telegram.WebhookMessage) error {
	text := `*🗺️ Mapa de Crimes*

Mapa com marcações dos crimes registrados na cidade no primeiro semestre de 2025.

- Link para o mapa: 
https://www.google.com/maps/d/viewer?mid=1Z-LoxrmX55O5_Odo1lRXoCcs5TOXifs

- Os dados criminais foram obtidos através do *Portal Transparência - Números sem Mistério* da SSP-SP. Link para o portal:
https://www.ssp.sp.gov.br/estatistica/consultas
	`

	return telegram.SendMessage(message.Chat.ID, text)
}

func CommandLastNews(message *telegram.WebhookMessage, modelPath string, limit int) error {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found at %s", modelPath)
	}

	log.Printf("INFO: Getting news...")
	allNews, err := news.GetLastNews()
	if err != nil {
		return fmt.Errorf("failed to get news from MEON: %v", err)
	}
	log.Printf("INFO: Fetched %d news items from MEON and Sampi", len(allNews))

	var filteredNews []news.News
	nlpService := nlp.NewNLPService(modelPath)
	for _, newsItem := range allNews {
		text := strings.Join([]string{newsItem.Origin, newsItem.Title, newsItem.Content}, " ")
		class, err := nlpService.ClassifyContent(text)
		if err != nil {
			return fmt.Errorf("failed to classify news '%s': %v", newsItem.Title, err)
		}

		switch class {
		case nlp.Good:
			filteredNews = append(filteredNews, newsItem)
		case nlp.Bad:
		default:
		}
	}

	text := "*📰 Últimas Notícias*\n\n"
	if len(filteredNews) == 0 {
		text += "*Nenhuma notícia encontrada.*"
	}

	for _, newsItem := range filteredNews {
		text += fmt.Sprintf(
			"- [%s](%s) %s\n",
			newsItem.Title,
			newsItem.Link,
			newsItem.Content,
		)
	}

	text += `\nFontes: [Meon](https://www.meon.com.br/noticias/rmvale), [Sampi](https://sampi.net.br/ovale/categoria/ultimas)`

	if limit > 0 && len(filteredNews) > limit {
		text += fmt.Sprintf("\n\n*Exibindo apenas as primeiras %d notícias filtradas.*", limit)
	}

	return telegram.SendMessage(message.Chat.ID, text)
}
