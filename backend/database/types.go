package database

import (
	"context"
	"time"
)

// Database is a common interface for read and write operations
type Database interface {
	// Write performs insert of new records
	Write(ctx context.Context, bucket string, data []Metrics)

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
