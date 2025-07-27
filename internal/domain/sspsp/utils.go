package sspsp

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
)

func GenerateCrimeStatisticsTable(data []CrimeStatistics) (output string) {
	output += "   Ano | Homicídios | Furtos | Roubos |  FRV  |\n"
	output += "-------+------------+--------+--------+-------|\n"

	// Ordena os dados por ano de forma decrescente
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data)-1; j++ {
			if data[j].Ano < data[j+1].Ano {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}

	for _, item := range data {
		output += fmt.Sprintf(
			" %5d | %10d | %6d | %6d | %5d |\n",
			item.Ano,
			item.Homicidio,
			item.Furto,
			item.Roubo,
			item.Frv,
		)
	}
	output += "-------+------------+--------+--------+-------|\n"
	output += "*  FRV: Furtos, Roubos de Veículos \n"
	return
}

func GenerateCrimeStatisticsImage(w, h int, data []CrimeStatistics) image.Image {
	canvas := gg.NewContext(w, h)

	canvas.SetColor(color.White)
	canvas.Clear()

	canvas.SetColor(color.Black)

	x := 10.0
	y := 20.0

	headers := []string{"Ano", "Homicídios", "Furtos", "Roubos", "FRV"}
	colW := float64(w / len(headers))

	ffHeader, _ := gg.LoadFontFace("./assets/OpenSans-VariableFont_wdth,wght.ttf", 16)
	canvas.SetFontFace(ffHeader)

	for _, header := range headers {
		canvas.DrawString(header, x, y)
		x += colW
	}

	ff, _ := gg.LoadFontFace("./assets/OpenSans-VariableFont_wdth,wght.ttf", 12)
	canvas.SetFontFace(ff)

	y += 20
	for _, item := range data {
		x = 10
		canvas.DrawString(fmt.Sprintf("%d", item.Ano), x, y)
		x += colW
		canvas.DrawString(fmt.Sprintf("%d", item.Homicidio), x, y)
		x += colW
		canvas.DrawString(fmt.Sprintf("%d", item.Furto), x, y)
		x += colW
		canvas.DrawString(fmt.Sprintf("%d", item.Roubo), x, y)
		x += colW
		canvas.DrawString(fmt.Sprintf("%d", item.Frv), x, y)
		y += 20
	}

	return canvas.Image()
}

func GenerateCrimeStatisticsDetailedTable(data []CrimeStatisticsDetailed) (output string) {
	output += " Jan | Fev | Mar | Abr | Mai | Jun | Jul | Ago | Set | Out | Nov | Dez | Delito \n"
	output += "-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+--------\n"

	// Ordena os dados por delito de forma crescente
	for j := 0; j < len(data)-1; j++ {
		if strings.Compare(data[j].Delito.Delito, data[j+1].Delito.Delito) > 0 {
			data[j], data[j+1] = data[j+1], data[j]
		}
	}

	for _, item := range data {
		output += fmt.Sprintf(
			" %3d | %3d | %3d | %3d | %3d | %3d | %3d | %3d | %3d | %3d | %3d | %3d | %s \n",
			item.Janeiro,
			item.Fevereiro,
			item.Marco,
			item.Abril,
			item.Maio,
			item.Junho,
			item.Julho,
			item.Agosto,
			item.Setembro,
			item.Outubro,
			item.Novembro,
			item.Dezembro,
			item.Delito.Delito,
		)
	}
	return
}
