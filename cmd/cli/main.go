package main

import (
	"flag"
	"fmt"
	"os"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/sspsp"
)

func main() {
	service := flag.String("service", "", "Service to be used: sspsp,cetesb")
	cityId := flag.Int("city_id", 0, "City ID to get the air quality")
	municipalityId := flag.Int("municipality_id", 0, "Municipality ID to get the police incidents")
	year := flag.Int("year", 0, "Year to get the data")
	detailed := flag.Bool("detailed", false, "Get detailed data")
	location := flag.Bool("location", false, "Get location data")
	flag.Parse()

	switch *service {
	case "cetesb":
		serviceCetesb(*cityId)
		os.Exit(0)

	case "sspsp":
		serviceSspsp(*year, *municipalityId, *detailed, *location)
		os.Exit(0)

	default:
		fmt.Printf("Usage: %s -service <sspsp|cetesb>\n", os.Args[0])
	}

	os.Exit(0)
}

func serviceCetesb(cityId int) {
	if cityId == 0 {
		fmt.Println("City ID is required")
		fmt.Printf("Usage: %s -service cetesb -city_id <city_id>", os.Args[0])
		os.Exit(1)
	}

	data, err := cetesb.GetQualarData(cityId)
	if err != nil {
		fmt.Printf("Error fetching: %v\n", err)
		os.Exit(1)
	}

	if len(data.Features) == 0 {
		fmt.Printf("Sem resultados para a cidade informada: %d\n", cityId)
		os.Exit(0)
	}

	fmt.Printf("Nome: %s\n", data.Features[0].Attributes.Nome)
	fmt.Printf("Indice qualidade do Ar: %f\n", data.Features[0].Attributes.Indice)
}

func serviceSspsp(year int, municipalityId int, detailed bool, location bool) {
	if municipalityId == 0 {
		fmt.Println("Municipality ID is required")
		fmt.Printf(
			"Usage: %s -service sspsp -municipality_id <municipality_id> [-detailed]",
			os.Args[0],
		)
		os.Exit(1)
	}

	if !detailed {
		data, err := sspsp.GetPoliceIncidentsCriminal(municipalityId)
		if err != nil {
			fmt.Printf("Error fetching: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%+v", data)
		return
	}

	if year == 0 {
		fmt.Println("Year is required when using detailed")
		fmt.Printf(
			"Usage: %s -service sspsp -municipality_id <municipality_id> -detailed -year <year>",
			os.Args[0],
		)
		os.Exit(1)
	}

	if location {
		data, err := sspsp.GetPoliceIncidentsByLocation(year)
		if err != nil {
			fmt.Printf("Error fetching: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%+v", data)
		return
	}

	data, err := sspsp.GetPoliceIncidentsCriminalDetailed(year, municipalityId)
	if err != nil {
		fmt.Println("Error fetching ")
		os.Exit(1)
	}

	fmt.Printf("%+v", data)
}
