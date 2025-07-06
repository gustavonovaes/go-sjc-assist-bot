package main

import (
	"strconv"
	"strings"
)

func ExtractCityIdFromMessage(message string) int {
	commandTokens := strings.Split(message, " ")

	var commandCityId int
	if len(commandTokens) < 2 {
		commandCityId = QUALAR_STATION_ID
	} else {
		commandCityId, _ = strconv.Atoi(commandTokens[1])
		if commandCityId == 0 {
			commandCityId = QUALAR_STATION_ID
		}
	}

	return commandCityId
}
