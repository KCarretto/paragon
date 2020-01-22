package cdn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/file"
)

const (
	MaxFileSize int64 = 10 << 20
	MaxMemSize        = 10 << 20
)

// HTTP handlers for file upload & download.
type HTTP struct {
	EntClient *ent.Client
}

// HandleFileUpload is an http.HandlerFunc which parses multipart forms and upserts file objects.
func (h HTTP) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)
	if err := r.ParseMultipartForm(MaxMemSize); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse multipart form: %v", err), http.StatusBadRequest)
		return
	}

	fileName := r.PostFormValue("fileName")
	if fileName == "" {
		http.Error(w, "must set valid value for 'fileName'", http.StatusBadRequest)
		return
	}

	fileQuery := h.EntClient.File.Query().Where(file.Name(fileName))
	exists := fileQuery.ExistX(ctx)

	f, _, err := r.FormFile("fileContent")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse file: %v", err), http.StatusBadRequest)
		return
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read file: %v", err), http.StatusBadRequest)
		return
	}

	var fileID int
	if exists {
		fileID = fileQuery.OnlyXID(ctx)
		h.EntClient.File.UpdateOneID(fileID).
			SetContent(content).
			SetSize(len(content)).
			SetLastModifiedTime(time.Now()).
			SaveX(ctx)
	} else {
		fileID = h.EntClient.File.Create().
			SetName(fileName).
			SetSize(len(content)).
			SetContent(content).
			SetLastModifiedTime(time.Now()).
			SaveX(ctx).ID
	}

	fmt.Fprintf(w, `{"data":{"file": {"id": %d}}}`, fileID)
}

// HandleFileDownload is an http.HandlerFunc which loads a file by name and serves it's content.
func (h HTTP) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filename := filepath.Base(r.URL.Path)
	if filename == "" || filename == "." || filename == "/" {
		http.Error(w, "invalid filename provided in request URI", http.StatusBadRequest)
		return
	}

	fileQuery := h.EntClient.File.Query().Where(file.Name(filename))

	if exists := fileQuery.ExistX(ctx); !exists {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	file := fileQuery.OnlyX(ctx)
	content := bytes.NewReader(file.Content)
	http.ServeContent(w, r, filename, file.LastModifiedTime, content)
}
