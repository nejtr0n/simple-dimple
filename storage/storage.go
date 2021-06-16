package storage

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	storageConsole = iota
	storageFile
)

const format = "source: %s, average: %d\r\n"

type Summary struct {
	Source string
	Value  int
}

type Storage interface {
	Store(summary Summary) error
	Close()
}

func NewStorage(storageType int) (Storage, error) {
	switch storageType {
	case storageConsole:
		return newConsoleStorage()
	case storageFile:
		return newFileStorage()
	default:
		return newConsoleStorage()
	}
}

func newConsoleStorage() (*consoleStorage, error) {
	return &consoleStorage{}, nil
}

type consoleStorage struct {
}

func (c *consoleStorage) Store(summary Summary) error {
	log.Printf(format, summary.Source, summary.Value)
	return nil
}

func (c consoleStorage) Close() {
}

func newFileStorage() (*fileStorage, error) {
	file, err := os.Create("output/data.txt")
	if err != nil {
		return nil, err
	}
	return &fileStorage{
		file: file,
	}, nil
}

type fileStorage struct {
	mu   sync.RWMutex
	file *os.File
}

func (f *fileStorage) Store(summary Summary) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	_, err := f.file.WriteString(fmt.Sprintf(format, summary.Source, summary.Value))
	if err != nil {
		return err
	}
	return nil
}

func (f *fileStorage) Close() {
	f.file.Close()
}
