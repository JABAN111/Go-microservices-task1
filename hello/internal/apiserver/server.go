package apiserver

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	mux    *http.ServeMux
	config *Config
}

func NewServer(config *Config) *Server {
	s := Server{
		http.NewServeMux(),
		config,
	}

	return &s
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, "pong\n")
		if err != nil {
			log.Fatalf("failed to write pong: %v", err)
			return
		}
	}
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		if len(name) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, err := fmt.Fprintf(w, "empty name\n")
			if err != nil {
				log.Fatalf("failed to write hello: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, "Hello, %s!\n", name)
		if err != nil {
			log.Fatalf("failed to write hello: %v", err)
		}

	}
}

func addRoutes(
	s *Server,
) {
	s.mux.HandleFunc("GET /ping", s.handlePing())
	s.mux.HandleFunc("GET /hello", s.handleHello())
}

func (s *Server) Run() {
	addRoutes(s)

	serverAddress := s.config.BindHost + ":" + s.config.BindPort
	log.Printf("Server started on address %s", serverAddress)

	err := http.ListenAndServe(serverAddress, s.mux)
	if err != nil {
		log.Panicf("Cannot started a server application, reason: %v", err)
	}
}
