package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/VladikAN/meteo-agent/config"
	"github.com/VladikAN/meteo-agent/database"

	"golang.org/x/crypto/acme/autocert"
)

var srv *http.Server
var db database.Database

// Start will start http server
func Start(cfg config.Settings) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warn("System interrupt or terminate signal")
		cancel()
	}()

	go func() {
		<-ctx.Done()
		shutdown()
	}()

	var err error

	// Configure database
	db, err = database.Start(cfg)
	if err != nil {
		log.Panicf("Failed to connect to the influx: %s", err)
	}
	defer db.Stop()

	// Configure and start service
	srv = &http.Server{Addr: cfg.Address, Handler: newRouter()}
	log.Infof("Server starting at %s", cfg.Address)

	if cfg.Ssl {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("autocert"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Whitelist...),
		}

		srv.TLSConfig = m.TLSConfig()
		log.Info("Setting up SSL for the whitelist")

		// serve HTTP, which will redirect automatically to HTTPS
		go http.ListenAndServe(":http", m.HTTPHandler(nil))
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}

	log.Warnf("Server was terminated or failed to start, %s", err)
}

func shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Infof("Server shutdown error, %s", err)
	}
}
