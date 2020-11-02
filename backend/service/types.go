package service

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
