package sspsp

import (
	"fmt"
	"os"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
)

const SJC_MUNICIPALITY_ID = 560

func CommandCrimes(message telegram.WebhookMessage) error {
	data, err := GetPoliceIncidentsCriminal(SJC_MUNICIPALITY_ID)
	if err != nil {
		fmt.Printf("Error fetching: %v\n", err)
		os.Exit(1)
	}

	return telegram.SendMessage(message.Chat.ID, fmt.Sprintf(
		"<code>\n%s\n</code>",
		GenerateCrimeStatisticsTable(data[:10]),
	))
}
