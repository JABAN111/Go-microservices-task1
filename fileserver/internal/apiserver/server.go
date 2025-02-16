package apiserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"yadro.com/course/internal/storage"
)

const (
	maxMemory = 32 << 20
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
		if err := r.ParseMultipartForm(maxMemory); err != nil {
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
func (s *Server) handleUpdateFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			log.Fatalf("failed to parse multipart form: %v", err)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Fatalf("failed to parse file: %v", err)
		}

		s.storage.Update(file, header) //TODO process
		log.Printf("file uploaded successfully")
		fmt.Fprintf(w, "file uploaded successfully")
	}
}

func (s *Server) handleListFiles() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		files, err := s.storage.GetFilesAsString()

		if err != nil {
			//TODO
			return
		}

		if _, err := fmt.Fprintf(w, files); err != nil {
			log.Fatalf("failed to write response: %v", err)
			return
		}
	}
}

func (s *Server) handleGetFile() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		filename := r.PathValue("filename")
		file, err := s.storage.Get(filename)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}(file)
		if err != nil {
			//TODO
			w.WriteHeader(http.StatusNotFound)
			if _, err := fmt.Fprintf(w, "file not found"); err != nil {
				return
			}
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, file); err != nil {
			//TODO
			return
		}
	}
}
func (s *Server) handleDeleteFile() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := s.storage.Delete(r.PathValue("filename"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		if _, err := fmt.Fprintf(w, "file deleted"); err != nil {
			return
		}
	}
}

func addRoutes(
	s *Server,
) {
	s.mux.HandleFunc("POST /files", s.handleSaveFile())
	s.mux.HandleFunc("PUT /files/{filename}", s.handleUpdateFile())
	s.mux.HandleFunc("GET /files/{filename}", s.handleGetFile())
	s.mux.HandleFunc("GET /files", s.handleListFiles())
	s.mux.HandleFunc("DELETE /files/{filename}", s.handleDeleteFile())
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
