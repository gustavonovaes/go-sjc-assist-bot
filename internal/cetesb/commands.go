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
			"Nome: %s\nIndice qualidade do Ar: %.2f - %s\n\n%s\n",
			res.Features[0].Attributes.Nome,
			res.Features[0].Attributes.Indice,
			res.Features[0].Attributes.Qualidade,
			"		0-40 	- Boa, 				41-80 - Moderada, \n"+
				"	 81-120 - Ruim,  		 121-200 - Muito Ruim,\n"+
				" 	 >200 - Péssima\n"+
				"\n",
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
