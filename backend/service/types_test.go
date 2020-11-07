package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToDbType(t *testing.T) {
	input := Metrics{
		Token: "token",
		Name:  "name",
		Data: []Measure{
			{Offset: 5, Temperature: 1, Humidity: 1},
			{Offset: 10, Temperature: 10, Humidity: 10},
		},
	}
	now := time.Now()

	result := toDbType(input, now)

	assert.Len(t, result, 2)
	assert.Equal(t, input.Name, result[0].Name)
	assert.Equal(t, input.Data[0].Offset, int(now.Sub(result[0].Date).Seconds()))
	assert.Equal(t, input.Data[0].Temperature, result[0].Temperature)
	assert.Equal(t, input.Data[0].Humidity, result[0].Humidity)

	assert.Equal(t, input.Name, result[1].Name)
	assert.Equal(t, input.Data[1].Offset, int(now.Sub(result[1].Date).Seconds()))
	assert.Equal(t, input.Data[1].Temperature, result[1].Temperature)
	assert.Equal(t, input.Data[1].Humidity, result[1].Humidity)
}
