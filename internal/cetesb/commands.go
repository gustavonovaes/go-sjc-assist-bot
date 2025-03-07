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
		return fmt.Errorf(
			"failed to get data from CETESB: %v\n Telegram response result: %v",
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
