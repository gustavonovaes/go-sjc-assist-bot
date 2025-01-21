package cetesb

import (
	"fmt"
	"strconv"
	"strings"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
)

const SJC_QUALAR_STATION_ID = 49

func CommandQualidadeAr(message telegram.WebhookMessage) error {
	commandCityId := extractCityId(message.Text)
	res, err := GetQualarData(commandCityId)
	if err != nil {
		telegram.SendMessage(message.Chat.ID, "Erro ao obter dados da CETESB")
		return fmt.Errorf("failed to get data from CETESB: %v", err)
	}

	return telegram.SendMessage(
		message.Chat.ID,
		fmt.Sprintf(
			"Nome: %s\nIndice qualidade do Ar: %f",
			res.Features[0].Attributes.Nome,
			res.Features[0].Attributes.Indice,
		),
	)
}

func extractCityId(message string) int {
	commandTokens := strings.Split(message, " ")

	var commandCityId int
	if len(commandTokens) < 2 {
		commandCityId = SJC_QUALAR_STATION_ID
	} else {
		commandCityId, _ = strconv.Atoi(commandTokens[1])
		if commandCityId == 0 {
			commandCityId = SJC_QUALAR_STATION_ID
		}
	}

	return commandCityId
}
