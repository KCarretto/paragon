package assets

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/spf13/afero"
)

// Bundler is an interface for handling multi-file (de)serialization and creation of a http.FileSystem
type Bundler interface {
	Bundle(...http.File) error
}

// TarGZBundler is the concrete implementation of Bundle using Tar and GZip
type TarGZBundler struct {
	Buffer *bytes.Buffer
}

// Bundle is used to add multiple files into a tar bundle
func (tb *TarGZBundler) Bundle(files ...http.File) error {
	tb.Buffer = &bytes.Buffer{}
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
	gzBuffer := &bytes.Buffer{}
	gw := gzip.NewWriter(gzBuffer)
	gw.Write(tb.Buffer.Bytes())
	if err := gw.Close(); err != nil {
		return err
	}
	tb.Buffer = gzBuffer
	return nil
}

// FileSystem is used to convert a targz bundle into an in memory http.FileSystem
func (tb *TarGZBundler) FileSystem() (http.FileSystem, error) {
	gr, err := gzip.NewReader(tb.Buffer)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(gr)
	if err := gr.Close(); err != nil {
		return nil, err
	}
	tr := tar.NewReader(bytes.NewBuffer(b))
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
