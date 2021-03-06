package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	HTTPPort *string   `mapstructure:"http_port" validate:"required"`
	AppName  *string   `mapstructure:"app_name" validate:"required"`
	Env      *string   `mapstructure:"environment" validate:"required"`
	LogLevel *string   `mapstructure:"log_level" validate:"required"`
	AWS      *AWS      `mapstructure:"aws" validate:"required"`
	DynamoDB *DynamoDB `mapstructure:"dynamodb" validate:"required"`
}

type AWS struct {
	Region *string `mapstructure:"region" validate:"required"`
}

type DynamoDB struct {
	Endpoint *string `mapstructure:"endpoint" validate:"required"`
	Table    *string `mapstructure:"table" validate:"required"`
}

func LoadConfiguration(path string, v *validator.Validate) *Configuration {
	file := viper.New()
	file.SetConfigFile(path)
	var config Configuration

	if err := file.ReadInConfig(); err != nil {
		log.Fatal("failing reading file, error: ", err)
		os.Exit(1)
	}
	if err := file.Unmarshal(&config); err != nil {
		log.Fatal("failing unmarshal configuration, error: ", err)
		os.Exit(1)
	}
	if err := v.Struct(config); err != nil {
		log.Fatal("failing validation, error: ", err)
		os.Exit(1)
	}
	return &config
}
