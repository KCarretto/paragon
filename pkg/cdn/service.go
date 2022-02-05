package cdn

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/event"
	"github.com/kcarretto/paragon/ent/file"
	"github.com/kcarretto/paragon/ent/link"
	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/service"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

const maxMemSize = 10 << 20

// Service provides HTTP handlers for the CDN.
type Service struct {
	Log   *zap.Logger
	Graph *ent.Client
	Auth  service.Authenticator
}

// HTTP registers http handlers for the CDN.
func (svc *Service) HTTP(router *http.ServeMux) {
	go svc.initFiles()

	upload := &service.Endpoint{
		Log:           svc.Log.Named("upload"),
		Authenticator: svc.Auth,
		Authorizer:    auth.NewAuthorizer().IsActivated(),
		Handler:       service.HandlerFn(svc.HandleFileUpload),
	}
	download := &service.Endpoint{
		Log:           svc.Log.Named("download"),
		Authenticator: svc.Auth,
		Authorizer:    auth.NewAuthorizer().IsActivated(),
		Handler:       service.HandlerFn(svc.HandleFileDownload),
	}
	links := &service.Endpoint{
		Log:           svc.Log.Named("links"),
		Authenticator: svc.Auth,
		Handler:       service.HandlerFn(svc.HandleLink),
	}

	router.Handle("/cdn/upload/", upload)
	router.Handle("/cdn/download/", download)
	router.Handle("/l/", http.StripPrefix("/l", links))
}

// HandleFileUpload is an http.HandlerFunc which parses multipart forms and upserts file objects.
func (svc Service) HandleFileUpload(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	if err := r.ParseMultipartForm(maxMemSize); err != nil {
		return fmt.Errorf("failed to parse multipart form: %w", err)
	}

	fileName := r.PostFormValue("fileName")
	if fileName == "" {
		return fmt.Errorf("must set valid value for 'fileName'")
	}

	fileQuery := svc.Graph.File.Query().Where(file.Name(fileName))
	exists := fileQuery.Clone().ExistX(ctx)

	f, _, err := r.FormFile("fileContent")
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var fileID int
	digestBytes := sha3.Sum256(content)
	digest := base64.StdEncoding.EncodeToString(digestBytes[:])
	contentType := http.DetectContentType(content)
	if exists {
		fileID = fileQuery.OnlyIDX(ctx)
		svc.Graph.File.UpdateOneID(fileID).
			SetContent(content).
			SetHash(digest).
			SetContentType(contentType).
			SetSize(len(content)).
			SetLastModifiedTime(time.Now()).
			SaveX(ctx)
	} else {
		fileID = svc.Graph.File.Create().
			SetName(fileName).
			SetSize(len(content)).
			SetContent(content).
			SetHash(digest).
			SetContentType(contentType).
			SetLastModifiedTime(time.Now()).
			SaveX(ctx).ID
	}
	// if we fail to create event, we don't wish to panic
	svc.Graph.Event.Create().
		SetOwner(auth.GetUser(ctx)).
		SetFileID(fileID).
		SetKind(event.KindUPLOAD_FILE).
		Save(ctx)
	fmt.Fprintf(w, `{"data":{"file": {"id": %d}}}`, fileID)
	return nil
}

// HandleFileDownload is an http.HandlerFunc which loads a file by name and serves it's content.
func (svc Service) HandleFileDownload(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	filename := filepath.Base(r.URL.Path)
	if filename == "" || filename == "." || filename == "/" {
		return fmt.Errorf("invalid filename provided in request URI")
	}

	fileQuery := svc.Graph.File.Query().Where(file.Name(filename))

	if exists := fileQuery.Clone().ExistX(ctx); !exists {
		return fmt.Errorf("file not found")
	}

	// If hash was provided, check to see if the file has been updated
	// Note that http.ServeContent should handle this, but we want to avoid the expensive DB query
	// where possible.
	if hash := r.Header.Get("If-None-Match"); hash != "" {
		if exists := fileQuery.Clone().Where(file.Hash(hash)).ExistX(ctx); exists {
			http.Error(w, "file has not been modified", http.StatusNotModified)
			return nil
		}
	}

	file := fileQuery.OnlyX(ctx)
	content := bytes.NewReader(file.Content)

	// Set hash of file
	digestBytes := sha3.Sum256(file.Content)
	w.Header().Set("Etag", base64.StdEncoding.EncodeToString(digestBytes[:]))

	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, filename, file.LastModifiedTime, content)
	return nil
}

// HandleLink is an http.HandlerFunc which loads a file by its link and serves it's content.
func (svc Service) HandleLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	alias := filepath.Base(r.URL.Path)
	if alias == "" || alias == "." || alias == "/" {
		return fmt.Errorf("invalid alias provided in request URI")
	}

	linkQuery := svc.Graph.Link.Query().Where(link.Alias(alias))

	if exists := linkQuery.ExistX(ctx); !exists {
		return fmt.Errorf("alias not found")
	}

	link := linkQuery.OnlyX(ctx)
	if link.Clicks == 0 || (link.ExpirationTime.Before(time.Now()) && !link.ExpirationTime.IsZero()) {
		svc.Graph.Link.DeleteOneID(link.ID).ExecX(ctx)
		return fmt.Errorf("alias not found")
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
	return nil
}

// populate with files from the cdn directory if they don't yet exist.
func (svc Service) initFiles() {
	var root = "cdn/"

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Skip errors or directories
		if err != nil || info == nil || info.IsDir() {
			return nil
		}

		// Remove root directory from file name
		name, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		// Skip if it already exists
		ctx := context.Background()
		if exists, err := svc.Graph.File.Query().Where(file.Name(name)).Exist(ctx); err != nil || exists {
			return nil
		}

		// Read the file
		content, err := ioutil.ReadFile(path)
		if err != nil || content == nil {
			return err
		}

		// Calculate the hash and content type
		digestBytes := sha3.Sum256(content)
		digest := base64.StdEncoding.EncodeToString(digestBytes[:])
		contentType := http.DetectContentType(content)

		// Create the file
		svc.Graph.File.Create().
			SetName(name).
			SetSize(int(info.Size())).
			SetContentType(contentType).
			SetContent(content).
			SetHash(digest).
			SetLastModifiedTime(time.Now()).
			Save(ctx)

		return nil
	})
}
