package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig - Load Config
func LoadConfig(filename string) {
	viper.SetConfigType("yaml")

	viper.SetConfigName(filename)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	viper.SetDefault("mysql.max_idle_connections", 5)
	viper.SetDefault("mysql.max_connections", 25)
	viper.SetDefault("mysql.max_lifetime", time.Second*5)
}
