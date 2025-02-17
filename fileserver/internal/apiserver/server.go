package apiserver

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
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

func (s *Server) handleSaveFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			log.Printf("Failed to parse multipart form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to parse file: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				log.Printf("Failed to close file: %v", err)
				return
			}
		}(file)

		if err := s.storage.Save(file, header); err != nil {
			log.Printf("Failed to save file: %v", err)
			http.Error(w, "Conflict", http.StatusConflict)
			return
		}

		log.Printf("File: %s uploaded successfully", header.Filename)
		w.WriteHeader(http.StatusCreated)
		if _, err := fmt.Fprintln(w, header.Filename); err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}
	}
}

func (s *Server) handleUpdateFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			log.Printf("Failed to parse multipart form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to parse file: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if err := s.storage.Update(file, header); err != nil {
			log.Printf("Failed to update file: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Println("File updated successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "File updated successfully")
	}
}

func (s *Server) handleListFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := s.storage.GetFilesAsString()
		if err != nil {
			log.Printf("Failed to list files: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintln(w, files); err != nil {
			log.Printf("Failed to write response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleGetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		file, err := s.storage.Get(filename)
		if err != nil {
			log.Printf("File not found: %v", err)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		w.Header().Set("Content-Transfer-Encoding", "binary")

		if _, err := io.Copy(w, file); err != nil {
			log.Printf("Failed to send file: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleDeleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		err := s.storage.Delete(filename)
		if err != nil {
			log.Printf("File not found: %v", err)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintln(w, "File deleted"); err != nil {
			log.Printf("Failed to write response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func addRoutes(s *Server) {
	s.mux.HandleFunc("POST /files", s.handleSaveFile())
	s.mux.HandleFunc("PUT /files/{filename}", s.handleUpdateFile())
	s.mux.HandleFunc("GET /files/{filename}", s.handleGetFile())
	s.mux.HandleFunc("GET /files", s.handleListFiles())
	s.mux.HandleFunc("DELETE /files/{filename}", s.handleDeleteFile())
	s.mux.Handle("/", http.NotFoundHandler())

	log.Println("Finished registering routes...")
}

func (s *Server) Run() {
	addRoutes(s)

	serverAddress := s.config.BindHost + ":" + s.config.BindPort
	log.Printf("File server started on address %s", serverAddress)

	err := http.ListenAndServe(serverAddress, s.mux)
	if err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
