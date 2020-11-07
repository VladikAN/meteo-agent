package main

import (
	"github.com/VladikAN/meteo-agent/config"
	"github.com/VladikAN/meteo-agent/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.Read("")

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	service.Start(config)
}
