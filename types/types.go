package types

type OBUData struct {
	OBUID     int     `json:"obuid"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUId int     `json:"obuid"`
	Unix  int64   `json:"unix"`
}
