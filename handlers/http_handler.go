package handlers

import (
	"fmt"
	"net/http"

	"github.com/chasehampton/gom/models"
)

type HttpHandler struct {
	*BaseHandler
}

type HttpTarget interface {
	Target
	GetUrl() string
}

type HttpTargetImpl struct {
	Url        string
	TargetPath string
}

func (ht HttpTargetImpl) GetTarget() string {
	return ht.TargetPath
}

func (ht HttpTargetImpl) GetUrl() string {
	return ht.Url
}

func (h *HttpHandler) ListFiles(act models.Action) ([]interface{}, error) {
	return nil, fmt.Errorf("Listing files from Http is not supported")
}

func (h *HttpHandler) UploadFiles(act models.Action) error {
	return fmt.Errorf("Uploading to Http is not supported")
}

func (h *HttpHandler) UploadFile() error {
	return fmt.Errorf("Uploading to Http is not supported")
}

func (h *HttpHandler) DownloadFiles(act models.Action) error {
	return fmt.Errorf("Download Files feature not implemented in HttpHandler...")
}

func (h *HttpHandler) DownloadFile(target Target) error {
	response, err := http.Get(target.(HttpTarget).GetUrl())
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failure. Status: %v", response.Status)
	}
	if !isValidFileContentType(response.Header.Get("Content-Type")) {
		return fmt.Errorf("Invalid content type: %v", response.Header.Get("Content-Type"))
	}
	defer response.Body.Close()
	return h.BaseHandler.WriteToFile(response.Body, target.GetTarget())
}

func (h *HttpHandler) DeleteFile() error {
	return fmt.Errorf("Deleting from Http is not supported")
}

func isValidFileContentType(contentType string) bool {
	validContentTypes := []string{
		"application/octet-stream",
		"application/pdf",
		"application/zip",
		"text/plain",
		"image/jpeg",
		"image/png",
	}

	for _, validType := range validContentTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}
