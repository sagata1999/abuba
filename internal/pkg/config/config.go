package config

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Log struct {
	Level      string
	DSN        string
	Tags       map[string]string
	File       string
	MaxAge     int
	MaxSize    int
	MaxBackups int
}

func InitConfig(path string) {
	viper.SetDefault("log", Log{
		Level:      "error",
		DSN:        "",
		File:       ".transcoder.log",
		MaxAge:     14,
		MaxSize:    500,
		MaxBackups: 3,
	})

	viper.SetConfigFile(path)
	err := viper.SafeWriteConfig()
	if err != nil {
		log.Errorf("unable save default config: %v", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Warn("config file not found")
		} else {
			log.Errorf("unable read config: %v", err)
		}
	}
}
