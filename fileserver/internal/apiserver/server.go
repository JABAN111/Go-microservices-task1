package apiserver

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"yadro.com/course/internal/storage"
)

const maxMemory = 32 * 1024 * 1024

type Server struct {
	mux     *http.ServeMux
	config  *Config
	storage *storage.Storage
}

func NewServer(config *Config, storage *storage.Storage) *Server {
	s := &Server{
		mux:     http.NewServeMux(),
		config:  config,
		storage: storage,
	}
	s.addRoutes()
	return s
}

func (s *Server) writeResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if _, err := fmt.Fprintln(w, message); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func safeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Printf("Failed to close file: %v", err)
	}
}

func (s *Server) handleSaveFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer safeClose(file)

		if err := s.storage.Save(file, header); err != nil {
			http.Error(w, "Conflict", http.StatusConflict)
			return
		}

		s.writeResponse(w, http.StatusCreated, header.Filename)
	}
}

func (s *Server) handleUpdateFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer safeClose(file)

		if err := s.storage.Update(file, header); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		s.writeResponse(w, http.StatusOK, "File updated successfully")
	}
}

func (s *Server) handleListFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := s.storage.GetFilesAsString()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.writeResponse(w, http.StatusOK, files)
	}
}

func (s *Server) handleGetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		file, err := s.storage.Get(filename)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer safeClose(file)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		w.Header().Set("Content-Transfer-Encoding", "binary")

		if _, err := io.Copy(w, file); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}
}

func (s *Server) handleDeleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		if err := s.storage.Delete(filename); err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		s.writeResponse(w, http.StatusOK, "File deleted")
	}
}

func (s *Server) addRoutes() {
	s.mux.HandleFunc("POST /files", s.handleSaveFile())
	s.mux.HandleFunc("PUT /files/{filename}", s.handleUpdateFile())
	s.mux.HandleFunc("GET /files/{filename}", s.handleGetFile())
	s.mux.HandleFunc("GET /files", s.handleListFiles())
	s.mux.HandleFunc("DELETE /files/{filename}", s.handleDeleteFile())
}

func (s *Server) Run() {
	serverAddress := s.config.BindHost + ":" + s.config.BindPort
	log.Printf("File server started on address %s", serverAddress)

	if err := http.ListenAndServe(serverAddress, s.mux); err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
