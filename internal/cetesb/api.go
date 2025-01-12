// Package cetesb fornece funcionalidades para acessar dados da CETESB.
// CETESB é a Agência do Governo do Estado de São Paulo responsável pelo controle, fiscalização, monitoramento e licenciamento de atividades geradoras de poluição.

package cetesb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const CETESB_QUALAR_URL = "https://arcgis.cetesb.sp.gov.br/server/rest/services/QUALAR/CETESB_QUALAR/MapServer/6/query"

func GetQualarData(cityId int) (QualarResponse, error) {
	urlBase, err := url.Parse(CETESB_QUALAR_URL)
	if err != nil {
		return QualarResponse{}, err
	}
	urlBase.RawQuery = buildQueryParams(cityId)

	response, err := http.Post(urlBase.String(), "application/json", nil)
	if err != nil {
		return QualarResponse{}, err
	}
	defer response.Body.Close()

	var data QualarResponse
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return QualarResponse{}, err
	}

	return data, nil
}

func buildQueryParams(cityId int) string {
	q := url.Values{}
	q.Add("f", "json")
	q.Add("where", fmt.Sprintf("ID=%d", cityId))
	// q.Add("returnGeometry", "false")
	// q.Add("spatialRel", "esriSpatialRelIntersects")
	q.Add("outFields", "*")

	return q.Encode()
}
