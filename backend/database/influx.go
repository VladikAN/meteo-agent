package database

import (
	"context"
	"fmt"

	"github.com/VladikAN/meteo-agent/config"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

type databaseInflux struct {
	client influxdb2.Client
}

// Start creates influxdb client
func Start(st config.Settings) Database {
	client := influxdb2.NewClient(st.InfluxHost, fmt.Sprintf("%s:%s", st.InfluxUser, st.InfluxPassword))
	return &databaseInflux{client: client}
}

func (db *databaseInflux) Stop() {
	db.client.Close()
}

func (db *databaseInflux) Write(ctx context.Context, bucket string, data []Metrics) {
	org := "" /* org not used for 1.8 */

	bucketAPI := db.client.BucketsAPI()
	if b, err := bucketAPI.FindBucketByName(ctx, bucket); err == nil && b == nil {
		_, err = bucketAPI.CreateBucket(ctx, &domain.Bucket{Name: bucket})
	}

	writeAPI := db.client.WriteAPI(org, bucket)

	for _, item := range data {
		pt := influxdb2.NewPoint(
			"agent",
			map[string]string{"name": item.Name},
			map[string]interface{}{"t": item.Temperature, "h": item.Humidity},
			item.Date)

		writeAPI.WritePoint(pt)
	}

	writeAPI.Flush()
}
