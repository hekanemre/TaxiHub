package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port    string `mapstructure:"port"`
	MongoDB struct {
		Host   string `mapstructure:"host"`
		DBName string `mapstructure:"dbname"`
	} `mapstructure:"mongodb"`
	IdleTimeout    time.Duration `mapstructure:"idleTimeout"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
	NearbyDistance int           `mapstructure:"nearbyDistance"`
}

func Read() *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Viper config error: %v\n", err)
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	return &appConfig
}
