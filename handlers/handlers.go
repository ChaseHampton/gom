package handlers

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/chasehampton/gom/logger"
	"github.com/chasehampton/gom/models"
	"github.com/chasehampton/gom/vaultclient"
)

type Handler interface {
	ListFiles(models.Action) ([]interface{}, error)
	UploadFiles(models.Action) error
	DownloadFiles(models.Action) error
	DeleteFile() error
}

type Target interface {
	GetTarget() string
}

type BaseHandler struct {
	locks  map[string]*sync.Mutex
	mu     sync.Mutex
	Vault  *vaultclient.VaultClient
	Logger *logger.Logger
}

func GetBaseHandler(vc *vaultclient.VaultClient, logger *logger.Logger) *BaseHandler {
	return &BaseHandler{
		locks:  make(map[string]*sync.Mutex),
		Vault:  vc,
		Logger: logger,
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

func (b *BaseHandler) GetUploadFiles(path string) ([]string, error) {
	var fileList []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	return fileList, err
}

func (b *BaseHandler) OpenFile(filePath string) (*os.File, error) {
	flock := b.GetLock(filePath)
	flock.Lock()
	defer flock.Unlock()
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (b *BaseHandler) DeleteFile(file interface{}) error {
	switch f := file.(type) {
	case *os.File:
		flock := b.GetLock(f.Name())
		flock.Lock()
		defer flock.Unlock()
		return os.Remove(f.Name())
	case string:
		flock := b.GetLock(file.(string))
		flock.Lock()
		defer flock.Unlock()
		return os.Remove(f)
	default:
		return fmt.Errorf("Unexpected file type")
	}
}

func (b *BaseHandler) ArchiveFile(file interface{}, dest string) error {
	switch f := file.(type) {
	case *os.File:
		flock := b.GetLock(f.Name())
		flock.Lock()
		defer flock.Unlock()
		return os.Rename(f.Name(), dest)
	case string:
		flock := b.GetLock(file.(string))
		flock.Lock()
		defer flock.Unlock()
		return os.Rename(f, dest)
	default:
		return fmt.Errorf("Unexpected file type")
	}
}

func (b *BaseHandler) WriteToFile(reader io.Reader, filename string) error {
	flock := b.GetLock(filename)
	flock.Lock()
	defer flock.Unlock()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, reader)
	return err
}
