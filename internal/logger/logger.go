package logger

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/config"
)

func CreateLogger(conf *config.Configuration) *log.Logger {
	logLevel, err := log.ParseLevel(*conf.LogLevel)
	if err != nil {
		log.WithField("logLevel", *conf.LogLevel).
			WithError(err).
			Fatal("invalid log level")
	}
	l := log.New()
	l.SetLevel(logLevel)
	l.Out = os.Stdout
	return l
}
