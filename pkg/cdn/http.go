package cdn

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/file"
	"github.com/kcarretto/paragon/ent/link"
	"golang.org/x/crypto/sha3"
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
	digestBytes := sha3.Sum256(content)
	digest := base64.StdEncoding.EncodeToString(digestBytes[:])
	contentType := http.DetectContentType(content)
	if exists {
		fileID = fileQuery.OnlyXID(ctx)
		h.EntClient.File.UpdateOneID(fileID).
			SetContent(content).
			SetHash(digest).
			SetContentType(contentType).
			SetSize(len(content)).
			SetLastModifiedTime(time.Now()).
			SaveX(ctx)
	} else {
		fileID = h.EntClient.File.Create().
			SetName(fileName).
			SetSize(len(content)).
			SetContent(content).
			SetHash(digest).
			SetContentType(contentType).
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
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, filename, file.LastModifiedTime, content)
}

// HandleLink is an http.HandlerFunc which loads a file by its link and serves it's content.
func (h HTTP) HandleLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	alias := filepath.Base(r.URL.Path)
	if alias == "" || alias == "." || alias == "/" {
		http.Error(w, "invalid alias provided in request URI", http.StatusBadRequest)
		return
	}

	linkQuery := h.EntClient.Link.Query().Where(link.Alias(alias))

	if exists := linkQuery.ExistX(ctx); !exists {
		log.Printf("alias: %v", alias)
		http.Error(w, "alias not found", http.StatusNotFound)
		return
	}

	link := linkQuery.OnlyX(ctx)
	if link.Clicks == 0 || (link.ExpirationTime.Before(time.Now()) && !link.ExpirationTime.IsZero()) {
		h.EntClient.Link.DeleteOneID(link.ID).ExecX(ctx)
		log.Printf("alias 2: %v", alias)
		log.Printf("alias 2: %v", time.Now())
		http.Error(w, "alias not found", http.StatusNotFound)
		return
	}
	// a click has been used!
	if link.Clicks > 0 {
		link.Update().
			AddClicks(-1).
			SaveX(ctx)
	}
	file := link.QueryFile().OnlyX(ctx)
	content := bytes.NewReader(file.Content)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeContent(w, r, file.Name, file.LastModifiedTime, content)
}
