package handlers

import (
	"fmt"
	"io/fs"
	"os"
	"sync"
)

type Handler interface {
	UploadFile() error
	DownloadFile() error
	DeleteFile() error
	ArchiveFile() error
}

type BaseHandler struct {
	locks map[string]*sync.Mutex
	mu    sync.Mutex
}

func GetBaseHandler() *BaseHandler {
	return &BaseHandler{
		locks: make(map[string]*sync.Mutex),
	}
}

func (b *BaseHandler) GetLock(path string) *sync.Mutex {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.locks[path]; !ok {
		b.locks[path] = &sync.Mutex{}
	}
	return b.locks[path]
}

func (b *BaseHandler) GetUploadFiles(path string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return files, err
	}
	return files, nil
}

func (b *BaseHandler) OpenFile(filename string) (*os.File, error) {
	flock := b.GetLock(filename)
	flock.Lock()
	defer flock.Unlock()
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (b *BaseHandler) DeleteFile(file interface{}) error {
	flock := b.GetLock(file.(string))
	flock.Lock()
	defer flock.Unlock()
	switch f := file.(type) {
	case *os.File:
		return os.Remove(f.Name())
	case string:
		return os.Remove(f)
	default:
		return fmt.Errorf("Unexpected file type")
	}
}
