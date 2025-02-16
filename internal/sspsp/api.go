// Package sspsp fornece funcionalidades para acessar dados da Secretaria de Segurança Pública do estado de São Paulo.
// fonte https://www.ssp.sp.gov.br/estatistica/consultas

package sspsp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/plandem/xlsx"
)

const BASE_URL = "https://www.ssp.sp.gov.br/v1"

var xlsCache = make(map[string][]byte)

func GetPoliceIncidentsCriminal(idMunicipality EnumMunicipality) ([]CrimeStatistics, error) {
	url := fmt.Sprintf(
		"%s/OcorrenciasAnuais/recuperaDadosMunicipio?idMunicipio=%d",
		BASE_URL,
		idMunicipality,
	)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var data GetPoliceIncidentsCriminalResponse
	json.NewDecoder(response.Body).Decode(&data)

	if !data.Success {
		return nil, fmt.Errorf("error fetching data")
	}

	return data.Data, nil
}

func GetPoliceIncidentsCriminalDetailed(
	year int,
	idMunicipality EnumMunicipality,
) ([]CrimeStatisticsDetailed, error) {
	url := fmt.Sprintf(
		"%s/OcorrenciasMensais/RecuperaDadosMensaisAgrupados?ano=%d&grupoDelito=6&tipoGrupo=MUNIC%%C3%%8DPIO&idGrupo=%d",
		BASE_URL,
		year,
		idMunicipality,
	)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var data GetPoliceIncidentsCriminalDetailedResponse
	json.NewDecoder(response.Body).Decode(&data)

	if !data.Success {
		return nil, fmt.Errorf("error fetching data")
	}

	if len(data.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	return data.Data[0].ListaDados, nil
}

func GetPoliceIncidentsByLocation(year int) ([]CrimeStatisticsByLocation, error) {
	//  https://www.ssp.sp.gov.br/assets/estatistica/transparencia/spDados/SPDadosCriminais_2024.xlsx

	// url := fmt.Sprintf(
	// 	"https://www.ssp.sp.gov.br/assets/estatistica/transparencia/spDados/SPDadosCriminais_%d.xlsx",
	// 	year,
	// )

	currentDirectory, _ := os.Getwd()

	filePath := fmt.Sprintf("%s/SPDadosCriminais_2024.xlsx", currentDirectory)
	fmt.Println("xlsx:", filePath)

	file, _ := os.OpenFile(filePath, os.O_RDONLY, 0644)

	// response, err := http.Get(url)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	xl, err := xlsx.Open(file)
	if err != nil {
		return nil, err
	}
	defer xl.Close()

	showMemoryUsage()

	for sheets := xl.Sheets(); sheets.HasNext(); {
		_, sheet := sheets.Next()
		defer sheet.Close()

		totalCols, totalRows := sheet.Dimension()

		fmt.Printf("sheet:%s totalRows:%d totalCols:%d\n", sheet.Name(), totalRows, totalCols)

		for rowIndex := 0; rowIndex < totalRows; rowIndex++ {
			for colIndex := 0; colIndex < totalCols; colIndex++ {
				// cell := sheet.Cell(rowIndex, colIndex)
				// fmt.Println(cell.String())
			}
		}

		showMemoryUsage()

	}

	return nil, nil
}

func showMemoryUsage() {
	currentMemoryUsage := &runtime.MemStats{}

	runtime.ReadMemStats(currentMemoryUsage)

	log.Printf(
		"Memory Usage: %v MB\n",
		currentMemoryUsage.Alloc/1024/1024,
	)
}
