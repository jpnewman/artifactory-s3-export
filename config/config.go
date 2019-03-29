package config

import (
	"fmt"

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
}
