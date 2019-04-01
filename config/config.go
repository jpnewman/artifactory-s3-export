package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func LoadConfig(filename string) {
	viper.SetConfigType("yaml")

	viper.SetConfigName(filename)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	viper.SetDefault("mysql.max_connections", 150)
	viper.SetDefault("mysql.max_lifetime", time.Second*5)
}
