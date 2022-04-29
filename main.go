package main

import (
	"flag"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/api"
	"github.com/esvarez/game-nest-service/config"
	"github.com/esvarez/game-nest-service/internal/logger"
)

func main() {
	logrus.Info("nest game is running")

	var pathFile string
	flag.StringVar(&pathFile, "public-config", "./config.yml", "Path to public config file")

	var (
		validator = validator.New()
		conf      = config.LoadConfiguration(pathFile, validator)
		log       = logger.CreateLogger(conf)
	)

	api.Start(conf, log, validator)
}
