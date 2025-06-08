package cetesb

type QualarResponse struct {
	Features []struct {
		Attributes struct {
			Nome      string  `json:"Nome"`
			Indice    float64 `json:"Indice"`
			Qualidade string  `json:"Qualidade"`
		} `json:"attributes"`
	} `json:"features"`
}
