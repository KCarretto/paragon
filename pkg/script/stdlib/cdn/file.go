package cdn

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/kcarretto/paragon/pkg/cdn"
)

// File struct is a struct that structs shit up.
type File struct {
	Path string

	cdn.Uploader
	cdn.Downloader

	content io.Reader
}

// Write uploads the file content to the CDN.
func (f *File) Write(p []byte) (int, error) {
	if err := f.Upload(f.Path, bytes.NewBuffer(p)); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Read downloads the file content from the CDN.
func (f *File) Read(p []byte) (int, error) {
	if f.content == nil {
		content, err := f.Download(f.Path)
		if err != nil {
			return 0, err
		}
		f.content = content
	}
	return f.content.Read(p)
}

// Name provides the name that was provided to openFile.
func (f *File) Name() string {
	return path.Base(f.Path)
}

// Chmod does not work.
func (f *File) Chmod(os.FileMode) error {
	return fmt.Errorf("CDN does not support chmod")
}

// Chown does not work.
func (f *File) Chown(uid, gid int) error {
	return fmt.Errorf("CDN does not support chown")
}

// Stat does not work.
func (f *File) Stat() (os.FileInfo, error) {
	return nil, fmt.Errorf("CDN does not support stat")
}

// Move does not work.
func (f *File) Move(dstPath string) error {
	return fmt.Errorf("CDN does not support move")
}

// Remove does not work.
func (f *File) Remove() error {
	return fmt.Errorf("CDN does not support remove")
}

// Sync does nothing.
func (f *File) Sync() error {
	return nil
}
