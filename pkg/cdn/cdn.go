package cdn

import (
	"io"
)

// An Uploader saves files by name. It is responsible for closing the file after writing if the
// file implements io.Closer.
type Uploader interface {
	Upload(name string, file io.Reader) error
}

// A Downloader retrieves a file by name.
type Downloader interface {
	Download(name string) (io.Reader, error)
}
