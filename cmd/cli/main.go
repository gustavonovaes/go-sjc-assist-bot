package main

import (
	"context"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/news"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/nlp"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/domain/sspsp"
)

const (
	cityId         = 49  // São José dos Campos
	municipalityId = 560 // São José dos Campos
	timeoutDefault = 5 * time.Second
	modelPath      = "model.gob"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("No command provided. %s", usage())
	}

	command := args[0]

	switch command {
	case "cetesb":
		data, err := cetesb.GetQualarData(cityId)
		if err != nil {
			log.Fatalf("Error fetching: %v\n", err)
		}

		if len(data.Features) == 0 {
			log.Fatalf("Sem resultados para a cidade informada: %d\n", cityId)
		}

		currentIndex := data.Features[0].Attributes.Indice
		currentQuality := data.Features[0].Attributes.Qualidade
		fmt.Println(strings.Repeat("=", 55))
		fmt.Printf("| %-15s | %-15s | %-15s |\n", "N1 - BOA", "N2 - MODERADA", "N3 - RUIM")
		fmt.Printf("| %-15s | %-15s | %-15s |\n", "N4 - MUITO RUIM", "N5 - PÉSSIMA", "")
		fmt.Println(strings.Repeat("=", 55))
		fmt.Printf("Índice atual: %s (%.2f)\n", currentQuality, currentIndex)
		fmt.Println(strings.Repeat("=", 55))

	case "sspsp":
		data, err := sspsp.GetPoliceIncidentsCriminal(municipalityId)
		if err != nil {
			log.Fatalf("Error fetching: %v\n", err)
		}

		tableStringTop10 := sspsp.GenerateCrimeStatisticsTable(data[:10])
		fmt.Println(tableStringTop10)
	case "sspsp:detailed":
		currentYear := time.Now().Year()
		if len(args) >= 2 {
			if year, err := strconv.Atoi(args[1]); err == nil {
				currentYear = year
			}
		}

		data, err := sspsp.GetPoliceIncidentsCriminalDetailed(currentYear, municipalityId)
		if err != nil {
			log.Fatalf("Error fetching detailed data: %v\n", err)
		}

		fmt.Println(sspsp.GenerateCrimeStatisticsDetailedTable(data))

	case "sspsp:image":
		filePath := "crime_statistics.png"
		if len(args) >= 2 {
			filePath = args[1]
		}

		data, err := sspsp.GetPoliceIncidentsCriminal(municipalityId)
		if err != nil {
			log.Fatalf("Error fetching detailed data: %v\n", err)
		}

		img := sspsp.GenerateCrimeStatisticsImage(800, 600, data)

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Error creating image file: %v\n", err)
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			log.Fatalf("Error encoding image: %v\n", err)
		}

	case "model:train":
		goodSubjects := []string{
			"emergencial", "povo", "apoio", "conscientização", "projeto", "infraestrutura", "conquista",
			"governo", "prefeitura", "municipio",
			"sjc", "são josé dos campos", "são josé", "estado de sp", "sp registra", "paulista",
			"investimento", "economia", "atinge meta", "operação acontece", "inaugura",
			"festa do", "fim de semana", "feira",
		}
		badSubjects := []string{
			"acidente", "violencia", "mort familia corpo", "assassinato", "roubo", "furt", "incendio", "atropel", "apreensão", "cachorr", "sexual", "mutilad", "agredid", "sem vida", "asfixiou", "em coma", "desaparecid",
			"quadrilha trafico", "confusão", "agressão", "polícia suspeita", "drogas", "armas", "tiroteio", "assalto", "sequestro", "explosão", "criminoso", "bolsonaro", "lula",
			"taubaté", "jacareí", "rio de janeiro",
			"tecnico de", "copa américa", "brasileirão", "copa de clubes", "fifa",
		}

		fmt.Printf("Training model with:\n - good subjects: %v\n - bad subjects: %v\n", goodSubjects, badSubjects)

		nlpService := nlp.NewNLPService(modelPath)
		nlpService.TrainModel(goodSubjects, badSubjects)

	case "model:test":
		text := strings.Join(args[1:], " ")
		fmt.Printf("Classifying content: '%s'\n", text)

		nlpService := nlp.NewNLPService(modelPath)
		class, err := nlpService.ClassifyContent(text)
		if err != nil {
			log.Fatalf("Error classifying content '%s': %v", text, err)
		}

		if strings.TrimSpace(text) == "" {
			log.Fatalf("Text to classify cannot be empty")
		}

		switch class {
		case nlp.Good:
			fmt.Printf("'%s' classified as 'Good'\n", text)
		case nlp.Bad:
			fmt.Printf("'%s' classified as 'Bad'\n", text)
		default:
			fmt.Printf("'%s' classified as 'Unknown'\n", text)
		}

	case "news":
		type newsResult struct {
			origin string
			news   []news.News
			err    error
		}

		sources := []struct {
			name string
			fn   func() ([]news.News, error)
		}{
			{"meon", news.GetMeonNews},
			{"sampi", news.GetSampiNews},
		}

		resultCh := make(chan newsResult, len(sources))

		ctx, cancel := context.WithTimeout(context.Background(), timeoutDefault)
		defer cancel()

		for _, source := range sources {
			go func(origin string, fn func() ([]news.News, error)) {
				defer func() {
					if r := recover(); r != nil {
						resultCh <- newsResult{origin: origin, err: fmt.Errorf("panic: %v", r)}
					}
				}()

				news, err := fn()
				select {
				case resultCh <- newsResult{origin: origin, news: news, err: err}:
				case <-ctx.Done():
				}
			}(source.name, source.fn)
		}

		for _, source := range sources {
			select {
			case <-ctx.Done():
				fmt.Printf("Timeout while fetching news for source '%s'\n", source.name)
			case result := <-resultCh:
				if result.err != nil {
					log.Fatalf("Error fetching news from %s: %v\n", result.origin, result.err)
					return
				}

				for _, newsItem := range result.news {
					fmt.Printf("[%s] - %s '%s'\n%s\n", newsItem.Origin, newsItem.Title, newsItem.Content, newsItem.Link)
				}
			}
		}
	case "news:filtered":
		var newsList []news.News
		meonNews, err := news.GetMeonNews()
		if err != nil {
			log.Fatalf("Error fetching Meon news: %v\n", err)
		}
		newsList = append(newsList, meonNews...)

		sampiNews, err := news.GetSampiNews()
		if err != nil {
			log.Fatalf("Error fetching Sampi news: %v\n", err)
		}
		newsList = append(newsList, sampiNews...)

		nlpService := nlp.NewNLPService(modelPath)

		for _, newsItem := range newsList {
			text := strings.Join([]string{newsItem.Origin, newsItem.Title, newsItem.Content}, " ")
			class, err := nlpService.ClassifyContent(text)
			if err != nil {
				fmt.Printf("   Error classifying news '%s': %v\n", newsItem.Title, err)
				continue
			}

			switch class {
			case nlp.Good:
				fmt.Printf("[ Good    ] %s\n", newsItem.Title)
			case nlp.Bad:
				fmt.Printf("[ Bad     ] %s\n", newsItem.Title)
			default:
				fmt.Printf("[ Unknown ] %s\n", newsItem.Title)
			}
		}

	default:
		fmt.Printf("No command provided. %s", usage())
	}

	os.Exit(0)
}

func usage() string {
	return fmt.Sprintf("Usage: %s <sspsp|sspsp:table|sspsp:detailed [year]|cetesb|news|news:filtered|model:train|model:test [text]>", os.Args[0])
}
