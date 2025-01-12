package cetesb

type QualarResponse struct {
	Features []struct {
		Attributes struct {
			Nome   string  `json:"Nome"`
			Indice float64 `json:"Indice"`
		} `json:"attributes"`
	} `json:"features"`
}
