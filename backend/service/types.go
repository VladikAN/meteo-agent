package service

import (
	"time"

	"github.com/VladikAN/meteo-agent/database"
)

// Metrics holds inbound parent data from meteo agent
type Metrics struct {
	Token string    `json:"token"`
	Name  string    `json:"name"`
	Data  []Measure `json:"data"`
}

// Measure holds measurements comming from meteo agent
type Measure struct {
	Offset      int     `json:"o"`
	Temperature float32 `json:"t"`
	Humidity    int     `json:"h"`
}

func toDbType(mt Metrics, start time.Time) []database.Metrics {
	var rst []database.Metrics
	for _, item := range mt.Data {
		rst = append(rst, database.Metrics{
			Name:        mt.Name,
			Date:        start.Add(time.Second * time.Duration(item.Offset) * -1),
			Temperature: item.Temperature,
			Humidity:    item.Humidity,
		})
	}

	return rst
}
