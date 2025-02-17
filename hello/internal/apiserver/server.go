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
	return &Server{
		mux:    http.NewServeMux(),
		config: config,
	}
}

func writeResponse(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)
	if _, err := fmt.Fprintf(w, format, args...); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, http.StatusOK, "pong\n")
	}
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		if name == "" {
			writeResponse(w, http.StatusBadRequest, "empty name\n")
			return
		}

		writeResponse(w, http.StatusOK, "Hello, %s!\n", name)
	}
}

func (s *Server) addRoutes() {
	s.mux.HandleFunc("GET /ping", s.handlePing())
	s.mux.HandleFunc("GET /hello", s.handleHello())
}

func (s *Server) Run() {
	s.addRoutes()

	serverAddress := s.config.BindHost + ":" + s.config.BindPort
	log.Printf("Server started on %s", serverAddress)

	if err := http.ListenAndServe(serverAddress, s.mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
