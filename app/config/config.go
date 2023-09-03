package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() {
	// Normal config
	viper.AddConfigPath("./app/config/")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type

	// Read Config
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Database config
	viper.AddConfigPath("./app/database/")
	viper.SetConfigName("datasource.mongo")
	viper.SetConfigType("json")

	viper.MergeInConfig()
}
