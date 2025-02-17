package storage

import (
	"fmt"
	"io"
	"log"
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

func (s *Storage) Save(file *os.File, filename string) error {
	filePath := filepath.Join(s.path, filename)

	if _, err := os.Stat(filePath); err == nil {
		return &fileErr{filepath: filePath, msg: "file already exists"}
	}

	return saveFile(file, filePath)
}

func (s *Storage) Get(filename string) (*os.File, error) {
	filePath := filepath.Join(s.path, filename)
	return os.Open(filePath)
}

func (s *Storage) Update(file *os.File, filename string) error {
	filePath := filepath.Join(s.path, filename)
	if _, err := os.Stat(filePath); err != nil {
		return &fileErr{filepath: filePath, msg: "file does not exist"}
	}

	return saveFile(file, filePath)
}

func (s *Storage) Delete(filename string) error {
	filePath := filepath.Join(s.path, filename)
	return os.Remove(filePath)
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
	filesList, err := s.GetFiles()
	if err != nil {
		return "", err
	}
	return strings.Join(filesList, "\n"), nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal("Failed to close file")
	}
}

// saveFile saves file with given name OR overrides it
func saveFile(file *os.File, filePath string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer closeFile(outFile)

	_, err = io.Copy(outFile, file)
	return err
}
