package main

import (
	"github.com/VladikAN/meteo-agent/config"
	"github.com/VladikAN/meteo-agent/service"
)

func main() {
	config := config.Read("")
	service.Start(config)
}
