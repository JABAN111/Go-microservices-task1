package storage

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Storage struct {
	path string
}

func NewStorage(path string) (*Storage, error) {
	if err := os.MkdirAll(path, 0750); err != nil {
		return nil, err
	}

	return &Storage{
		path: path,
	}, nil
}

// TODO
func (s *Storage) Save(file multipart.File, header *multipart.FileHeader) error {
	filePath := filepath.Join(s.path, header.Filename)
	outFile, err := os.Create(filePath)
	defer closeFile(outFile)
	if err != nil {
		return err
	}
	if _, err := io.Copy(outFile, file); err != nil {
		return err
	}

	return nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal("Failed to close file")
	}
}
