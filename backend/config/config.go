package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Settings holds all application settings
type Settings struct {
	Address   string
	Ssl       bool
	Whitelist []string
}

// Read is called first to read all settings
func Read(path string) (cfg Settings) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path)

	viper.SetEnvPrefix("us")
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
