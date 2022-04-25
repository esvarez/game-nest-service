package main

import (
	"flag"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/internal/config"
)

func main() {
	log.Info("nest game is running")

	var pathFile string
	flag.StringVar(&pathFile, "public-file", "./config.yml", "path config to file")

	var (
		validator = validator.New()
		conf      = config.LoadConfiguration(pathFile, validator)
		_         = config.CreateLogger(conf)
	)
}
