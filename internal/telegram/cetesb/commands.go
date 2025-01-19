package cetesb

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
)

const SJC_QUALAR_STATION_ID = 49

func CommandQualidadeAr(message telegram.WebhookMessage) {
	commandToken := strings.Split(message.Text, " ")[1]
	commandCityId, _ := strconv.Atoi(commandToken)

	if commandCityId == 0 {
		commandCityId = SJC_QUALAR_STATION_ID
	}

	res, err := cetesb.GetQualarData(SJC_QUALAR_STATION_ID)
	if err != nil {
		telegram.SendMessage(message.Chat.ID, "Erro ao obter dados da CETESB")
		return
	}

	log.Printf("Res: %v", res)

	telegram.SendMessage(
		message.Chat.ID,
		fmt.Sprintf(
			"Nome: %s\nIndice qualidade do Ar: %f",
			res.Features[0].Attributes.Nome,
			res.Features[0].Attributes.Indice,
		),
	)
}
