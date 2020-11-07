package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/VladikAN/meteo-agent/database"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

var r *chi.Mux

type dbMock struct {
	createError error
	writeError  error
}

func (db *dbMock) Stop() {}

func (db *dbMock) CreateDatabaseIfMissed(bucket string) error {
	return db.createError
}

func (db *dbMock) Write(bucket string, data []database.Metrics) error {
	return db.writeError
}

func TestMain(m *testing.M) {
	r = chi.NewRouter()
	r.Post("/", postMetrics)

	code := m.Run()
	os.Exit(code)
}

func TestPostWithEmpty(t *testing.T) {
	rq, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	rq.Header.Set("content-type", "application/json")

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	assert.Equal(t, http.StatusBadRequest, rsp.Code)
	assert.Contains(t, rsp.Body.String(), "Input data is not valid")
}

func TestPostWithNoToken(t *testing.T) {
	rq, _ := http.NewRequest("POST", "/", strings.NewReader(`{}`))
	rq.Header.Set("content-type", "application/json")

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	assert.Equal(t, http.StatusBadRequest, rsp.Code)
	assert.Contains(t, rsp.Body.String(), "token value is empty")
}

func TestPostWithNoName(t *testing.T) {
	rq, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Token":"token"}`))
	rq.Header.Set("content-type", "application/json")

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	assert.Equal(t, http.StatusBadRequest, rsp.Code)
	assert.Contains(t, rsp.Body.String(), "name value is empty")
}

func TestPostWithNoData(t *testing.T) {
	rq, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Token":"token","Name":"name"}`))
	rq.Header.Set("content-type", "application/json")

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	assert.Equal(t, http.StatusOK, rsp.Code)
	assert.Contains(t, rsp.Body.String(), "No data received")
}

func TestPostDatabaseNotCreated(t *testing.T) {
	rq, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Token":"token","Name":"name","Data":[{"o":1,"t":2,"h":3}]}`))
	rq.Header.Set("content-type", "application/json")

	db = &dbMock{createError: errors.New("test-error")}

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	assert.Equal(t, http.StatusInternalServerError, rsp.Code)
	assert.Contains(t, rsp.Body.String(), "Error while writing data")
}

func TestPostWriteFailed(t *testing.T) {
	rq, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Token":"token","Name":"name","Data":[{"o":1,"t":2,"h":3}]}`))
	rq.Header.Set("content-type", "application/json")

	db = &dbMock{writeError: errors.New("test-error")}

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	assert.Equal(t, http.StatusInternalServerError, rsp.Code)
	assert.Contains(t, rsp.Body.String(), "Error while writing data")
}

func TestPostSaved(t *testing.T) {
	rq, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Token":"token","Name":"name","Data":[{"o":1,"t":2,"h":3}]}`))
	rq.Header.Set("content-type", "application/json")

	db = &dbMock{}

	rsp := httptest.NewRecorder()
	r.ServeHTTP(rsp, rq)

	assert.Equal(t, http.StatusOK, rsp.Code)
	assert.Contains(t, rsp.Body.String(), "Data is saved")
}
