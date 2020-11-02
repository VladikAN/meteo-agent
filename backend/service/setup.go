package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VladikAN/meteo-agent/config"

	"golang.org/x/crypto/acme/autocert"
)

var srv *http.Server

// Start will start http server
func Start(cfg config.Settings) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("WARN System interrupt or terminate signal")
		cancel()
	}()

	go func() {
		<-ctx.Done()
		shutdown()
	}()

	// Configure and start service
	srv = &http.Server{Addr: cfg.Address, Handler: newRouter()}
	log.Printf("Server starting at %s", cfg.Address)

	var err error
	if cfg.Ssl {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("autocert"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Whitelist...),
		}

		srv.TLSConfig = m.TLSConfig()
		log.Printf("Setting up SSL for the whitelist")

		// serve HTTP, which will redirect automatically to HTTPS
		go http.ListenAndServe(":http", m.HTTPHandler(nil))
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}

	log.Printf("Server was terminated or failed to start, %s", err)
}

func shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error, %s", err)
	}
}
