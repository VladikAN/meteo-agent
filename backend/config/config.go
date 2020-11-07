package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Settings holds all application settings
type Settings struct {
	// Debug will enable debug level log messages
	Debug bool

	// Address is an endpoint for the inbound connections
	Address string

	// Ssl is to use ssl or not
	Ssl bool

	// Whitelist is for lets encrypt domain names if Ssl is true
	Whitelist []string

	// InfluxHost is for target host with InfluxDB installed
	InfluxHost string

	// InfluxUser is for auth username to the InfluxDB
	InfluxUser string

	// INfluxPassword is for auth password for the InfluxDB
	InfluxPassword string
}

// Read is called first to read all settings
func Read(path string) (cfg Settings) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path)

	viper.SetEnvPrefix("ma")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error on reading config: %s", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("error on unmarshaling config: %s", err))
	}

	return cfg
}
