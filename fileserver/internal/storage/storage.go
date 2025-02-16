package storage

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Storage struct {
	path string
}

type fileErr struct {
	filepath string
	msg      string
}

func (e *fileErr) Error() string {
	return fmt.Sprintf("Error with file: %s, reason: %s", e.filepath, e.msg)
}

func NewStorage(path string) (*Storage, error) {
	if err := os.MkdirAll(path, 0750); err != nil {
		return nil, err
	}

	return &Storage{
		path: path,
	}, nil
}

// TODO process duplicate
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

func (s *Storage) Get(filename string) (*os.File, error) {
	filePath := filepath.Join(s.path, filename)
	outFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return outFile, nil
}

func (s *Storage) Update(file multipart.File, header *multipart.FileHeader) error {
	filePath := filepath.Join(s.path, header.Filename)
	if _, err := os.Stat(filePath); err != nil {
		return &fileErr{filepath: filePath, msg: "File are not exists"}
	}

	if err := s.Save(file, header); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(filename string) error {
	filePath := filepath.Join(s.path, filename)

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetFiles() ([]string, error) {
	files, err := os.ReadDir(s.path)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			res = append(res, f.Name())
		}
	}

	sort.Strings(res)

	return res, nil
}

func (s *Storage) GetFilesAsString() (string, error) {
	res, err := s.GetFiles()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, v := range res {
		sb.WriteString(v)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal("Failed to close file")
	}
}
