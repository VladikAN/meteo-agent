package database

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type databaseInflux struct {
	client client.Client
}

// Database is a common interface for read and write operations
type Database interface {
	// Write performs insert of new records
	Write(bucket string, data []Metrics) error

	// Stop will close db connection
	Stop()
}

// Metrics is a bucket data to store
type Metrics struct {
	Name        string
	Date        time.Time
	Temperature float32
	Humidity    int
}
