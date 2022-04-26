package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
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
			log.Println(err)
		}
	}()
	log.Printf("api is ready to handle requests: %s", s.Addr)
	s.gracefulShutdown()
}

func (s *server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Printf("api is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("api could not shutdown gracefully %s", err.Error())
	}
	log.Printf("api shutdown gracefully")
}
