package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Post("/", postMetrics)
	return r
}

func postMetrics(w http.ResponseWriter, r *http.Request) {
	var m Metrics

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w,
			fmt.Errorf("Error while reading input data: %s", err).Error(), http.StatusBadRequest)
		return
	}

	if result, reason := validatePostMetrics(m); result {
		http.Error(w, reason, http.StatusBadRequest)
		return
	}

	log.Printf("New message for '%s', '%s' agent, and %d measures", m.Token, m.Name, len(m.Data))

	if len(m.Data) == 0 {
		w.WriteHeader(http.StatusOK) // early exit
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(``))
}

func validatePostMetrics(m Metrics) (bool, string) {
	if len(m.Token) == 0 {
		return false, "token value is empty"
	}

	if len(m.Name) == 0 {
		return false, "name value is empty"
	}

	return true, ""
}