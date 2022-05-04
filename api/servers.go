package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	*http.Server
}

func newServer(listening string, mux *mux.Router) *server {
	return &server{
		Server: &http.Server{
			Addr:         ":" + listening,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

func (s *server) Start() {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("server error")
		}
	}()
	logrus.Infof("api is ready to handle requests: %s", s.Addr)
	s.gracefulShutdown()
}

func (s *server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	logrus.Infof("api is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatalf("api could not shutdown gracefully %s", err.Error())
	}
	logrus.Info("api shutdown gracefully")
}
