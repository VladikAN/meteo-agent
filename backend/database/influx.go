package database

import (
	"fmt"

	"github.com/VladikAN/meteo-agent/config"

	client "github.com/influxdata/influxdb1-client/v2"
)

// Start creates influxdb client
func Start(st config.Settings) (Database, error) {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     st.InfluxHost,
		Username: st.InfluxUser,
		Password: st.InfluxPassword,
	})

	if err != nil {
		return nil, err
	}

	return &databaseInflux{client: client}, nil
}

func (db *databaseInflux) Stop() {
	db.client.Close()
}

func (db *databaseInflux) CreateDatabaseIfMissed(bucket string) error {
	q := client.NewQuery(fmt.Sprintf(`CREATE DATABASE "%s"`, bucket), "", "")
	if resp, err := db.client.Query(q); err != nil {
		return err
	} else if len(resp.Err) != 0 {
		return fmt.Errorf(resp.Err)
	}

	q = client.NewQuery(fmt.Sprintf(`CREATE RETENTION POLICY "3month" ON "%s" DURATION 12w REPLICATION 1 DEFAULT`, bucket), bucket, "")
	if resp, err := db.client.Query(q); err != nil {
		return err
	} else if len(resp.Err) != 0 {
		return fmt.Errorf(resp.Err)
	}

	return nil
}

func (db *databaseInflux) Write(bucket string, data []Metrics) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision:       "s",
		RetentionPolicy: "3month",
		Database:        bucket,
	})

	if err != nil {
		return err
	}

	for _, item := range data {
		pt, err := client.NewPoint(
			"agent",
			map[string]string{"name": item.Name},
			map[string]interface{}{"t": item.Temperature, "h": item.Humidity},
			item.Date)

		if err != nil {
			return err
		}

		bp.AddPoint(pt)
	}

	return db.client.Write(bp)
}
