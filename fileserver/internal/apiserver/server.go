package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"yadro.com/course/internal/storage"
)

type Server struct {
	mux     *http.ServeMux
	config  *Config
	storage *storage.Storage
}

func NewServer(config *Config, storage *storage.Storage) *Server {
	s := Server{
		mux:     http.NewServeMux(),
		config:  config,
		storage: storage,
	}

	return &s
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, "pong\n")
		if err != nil {
			log.Fatalf("failed to write pong: %v", err)
		}
	}
}

func (s *Server) handleSaveFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		//TODO theoretically it's possible to delete, cause default max memory 32 mb too
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Fatalf("failed to parse multipart form: %v", err)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Fatalf("failed to parse file: %v", err)
		}

		s.storage.Save(file, header)
		log.Printf("file uploaded successfully")
		fmt.Fprintf(w, "pong\n")

	}
}

func addRoutes(
	s *Server,
) {
	s.mux.HandleFunc("POST /files", s.handleSaveFile())
	s.mux.HandleFunc("GET /ping", s.handlePing())
	s.mux.Handle("/", http.NotFoundHandler())
	log.Println("Finished registring routes...")

}

func (s *Server) Run() {
	addRoutes(s)

	serverAddress := s.config.BindHost + ":" + s.config.BindPort
	log.Printf("File server started on address %s", serverAddress)

	err := http.ListenAndServe(serverAddress, s.mux)
	if err != nil {
		log.Panicf("Cannot started a server application, reason: %v", err)
	}
}
