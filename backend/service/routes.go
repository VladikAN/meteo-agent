package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

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
		log.Warnf("Error while reading request data: %s", err)
		http.Error(w, fmt.Sprintf("Input data is not valid: %s", err), http.StatusBadRequest)
		return
	}

	if valid, message := validatePostMetrics(m); !valid {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	if len(m.Data) == 0 {
		log.Debugf("Empty data received for %s:%s", m.Token, m.Name)
		w.Write([]byte("No data received"))
		return
	}

	log.Infof("New message for '%s':'%s' agent with %d record(s)", m.Token, m.Name, len(m.Data))

	data := toDbType(m, time.Now())
	if err = db.Write(m.Token, data); err != nil {
		log.Errorf("Error while writing data: %s", err)
		http.Error(w, fmt.Sprintf("Error while writng data"), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Data is saved"))
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
