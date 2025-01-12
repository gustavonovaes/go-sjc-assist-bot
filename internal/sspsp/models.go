package sspsp

type EnumRegions int

const (
	AllRegions         EnumRegions = 0
	Capital            EnumRegions = 1
	GrandeSaoPaulo     EnumRegions = 2
	SaoJoseDosCampos   EnumRegions = 3
	Campinas           EnumRegions = 4
	RibeiraoPreto      EnumRegions = 5
	Bauru              EnumRegions = 6
	SaoJoseDoRioPreto  EnumRegions = 7
	Santos             EnumRegions = 8
	Sorocaba           EnumRegions = 9
	PresidentePrudente EnumRegions = 10
	Piracicaba         EnumRegions = 11
	Aracatuba          EnumRegions = 12
)

type EnumMunicipality int

const (
	AllMunicipalitys     EnumMunicipality = 0
	Aparecida            EnumMunicipality = 28
	Arapeí               EnumMunicipality = 36
	Areias               EnumMunicipality = 41
	Bananal              EnumMunicipality = 56
	Caçapava             EnumMunicipality = 97
	CachoeiraPaulista    EnumMunicipality = 98
	CamposdoJordão       EnumMunicipality = 111
	Canas                EnumMunicipality = 114
	Caraguatatuba        EnumMunicipality = 121
	Cruzeiro             EnumMunicipality = 150
	Cunha                EnumMunicipality = 152
	Guaratinguetá        EnumMunicipality = 211
	Igaratá              EnumMunicipality = 233
	Ilhabela             EnumMunicipality = 237
	Jacareí              EnumMunicipality = 280
	Jambeiro             EnumMunicipality = 285
	Lagoinha             EnumMunicipality = 300
	Lavrinhas            EnumMunicipality = 303
	Lorena               EnumMunicipality = 309
	MonteiroLobato       EnumMunicipality = 358
	NatividadedaSerra    EnumMunicipality = 365
	Paraibuna            EnumMunicipality = 403
	Pindamonhangaba      EnumMunicipality = 430
	Piquete              EnumMunicipality = 434
	Potim                EnumMunicipality = 458
	Queluz               EnumMunicipality = 473
	RedençãodaSerra      EnumMunicipality = 477
	Roseira              EnumMunicipality = 500
	SantaBranca          EnumMunicipality = 517
	SantoAntônioDoPinhal EnumMunicipality = 542
	SãoBentoDoSapucaí    EnumMunicipality = 546
	SãoJoséDoBarreiro    EnumMunicipality = 557
	SãoJoséDosCampos     EnumMunicipality = 560
	SãoLuísDoParaitinga  EnumMunicipality = 562
	SãoSebastião         EnumMunicipality = 569
	Silveiras            EnumMunicipality = 582
	Taubaté              EnumMunicipality = 607
	Tremembé             EnumMunicipality = 616
	Ubatuba              EnumMunicipality = 624
)

type CrimeStatistics struct {
	IDUnidade    EnumMunicipality `json:"idUnidade"`
	Ano          int              `json:"ano"`
	Homicidio    int              `json:"homicidio"`
	Furto        int              `json:"furto"`
	Roubo        int              `json:"roubo"`
	FurtoVeiculo int              `json:"furtoVeiculo"`
	RouboVeiculo int              `json:"rouboVeiculo"`
	Frv          int              `json:"frv"`
}

type GetPoliceIncidentsCriminalResponse struct {
	Success bool              `json:"success"`
	Data    []CrimeStatistics `json:"data"`
}

type CrimeStatisticsDetailed struct {
	IDOcorrenciaMensal int `json:"idOcorrenciaMensal"`
	IDDelito           int `json:"idDelito"`
	IDDistrito         int `json:"idDistrito"`
	Ano                int `json:"ano"`
	Janeiro            int `json:"janeiro"`
	Fevereiro          int `json:"fevereiro"`
	Marco              int `json:"marco"`
	Abril              int `json:"abril"`
	Maio               int `json:"maio"`
	Junho              int `json:"junho"`
	Julho              int `json:"julho"`
	Agosto             int `json:"agosto"`
	Setembro           int `json:"setembro"`
	Outubro            int `json:"outubro"`
	Novembro           int `json:"novembro"`
	Dezembro           int `json:"dezembro"`
	Publicado          int `json:"publicado"`
	Total              int `json:"total"`
	Delito             struct {
		IDDelito      int    `json:"idDelito"`
		IDGrupoDelito int    `json:"idGrupoDelito"`
		Delito        string `json:"delito"`
		Ordem         int    `json:"ordem"`
	} `json:"delito"`
}

type GetPoliceIncidentsCriminalDetailedResponse struct {
	Success bool `json:"success"`
	Data    []struct {
		Ano        int                       `json:"ano"`
		ListaDados []CrimeStatisticsDetailed `json:"listaDados"`
	} `json:"data"`
}

type CrimeStatisticsByLocation struct {
	NOMEDEPARTAMENTO            string `json:"NOME_DEPARTAMENTO"`
	NOMESECCIONAL               string `json:"NOME_SECCIONAL"`
	NOMEDELEGACIA               string `json:"NOME_DELEGACIA"`
	NOMEMUNICIPIO               string `json:"NOME_MUNICIPIO"`
	ANOBO                       int    `json:"ANO_BO"`
	NUMBO                       string `json:"NUM_BO"`
	DATAREGISTRO                string `json:"DATA_REGISTRO"`
	DATAOCORRENCIABO            string `json:"DATA_OCORRENCIA_BO"`
	HORAOCORRENCIABO            string `json:"HORA_OCORRENCIA_BO"`
	DESCPERIODO                 string `json:"DESC_PERIODO"`
	DESCRSUBTIPOLOCAL           string `json:"DESCR_SUBTIPOLOCAL"`
	BAIRRO                      string `json:"BAIRRO"`
	LOGRADOURO                  string `json:"LOGRADOURO"`
	NUMEROLOGRADOURO            int    `json:"NUMERO_LOGRADOURO"`
	LATITUDE                    string `json:"LATITUDE"`
	LONGITUDE                   string `json:"LONGITUDE"`
	NOMEDELEGACIACIRCUNSCRIO    string `json:"NOME_DELEGACIA_CIRCUNSCRIÇÃO"`
	NOMEDEPARTAMENTOCIRCUNSCRIO string `json:"NOME_DEPARTAMENTO_CIRCUNSCRIÇÃO"`
	NOMESECCIONALCIRCUNSCRIO    string `json:"NOME_SECCIONAL_CIRCUNSCRIÇÃO"`
	NOMEMUNICIPIOCIRCUNSCRIO    string `json:"NOME_MUNICIPIO_CIRCUNSCRIÇÃO"`
	RUBRICA                     string `json:"RUBRICA"`
	DESCRCONDUTA                string `json:"DESCR_CONDUTA"`
	NATUREZAAPURADA             string `json:"NATUREZA_APURADA"`
	MESESTATISTICA              int    `json:"MES_ESTATISTICA"`
	ANOESTATISTICA              int    `json:"ANO_ESTATISTICA"`
}
