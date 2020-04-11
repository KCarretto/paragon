package assets

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/spf13/afero"
)

// Bundle is an interface for handling multi-file (de)serialization and creation of a http.FileSystem
type Bundle interface {
	AddFiles(string, http.File) error
	Encode() string
	Decode(string) error
	ToFileSystem() (http.FileSystem, error)
}

// TarBundle is the concrete implementation of Bundle using Tar
type TarBundle struct {
	Buffer *bytes.Buffer
}

// AddFiles is used to add multiple files into a tar bundle
func (tb *TarBundle) AddFiles(files ...http.File) error {
	tw := tar.NewWriter(tb.Buffer)
	for _, file := range files {
		info, err := file.Stat()
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		hdr := &tar.Header{
			Name: info.Name(),
			Mode: 0644,
			Size: int64(len(body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(body); err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}

// Encode is used to encode a tar bundle into a single string
func (tb *TarBundle) Encode() string {
	return tb.Buffer.String()
}

// Decode is used to decode a tar bundle from a string
func (tb *TarBundle) Decode(buf string) error {
	tb.Buffer = bytes.NewBufferString(buf)
	return nil
}

// ToFileSystem is used to convert a tar bundle into an in memory http.FileSystem
func (tb *TarBundle) ToFileSystem() (http.FileSystem, error) {
	tr := tar.NewReader(tb.Buffer)
	fs := afero.NewMemMapFs()
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, err
		}
		f, err := fs.Create(hdr.Name)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(f, tr); err != nil {
			return nil, err
		}
	}
	return afero.NewHttpFs(fs), nil
}
