package handler

type Uni struct {
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	Isocode   string            `json:"isocode"`
	Webpages  []string          `json:"webpages"`
	Languages map[string]string `json:"languages"`
	Map       string            `json:"map"`
}

type Diag struct {
	UniApi     string  `json:"universitiesapi"`
	CountryApi string  `json:"countriesapi"`
	Version    string  `json:"version"`
	Uptime     float64 `json:"uptime"`
}
