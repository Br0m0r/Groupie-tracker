package models

type Coordinates struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Address string  `json:"address"`
}

type NominatimResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}
