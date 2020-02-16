package ssh

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/sftp"
)

type File struct {
	*sftp.File

	session *sftp.Client
}

// Move a file that is open via SFTP.
func (f *File) Move(dstPath string) error {
	dir := path.Dir(dstPath)

	if err := f.session.MkdirAll(dstPath); err != nil {
		return fmt.Errorf("failed to create parent directory %q: %w", dir, err)
	}

	err := f.session.Rename(f.Name(), dstPath)
	if err != nil {
		return err
	}

	err = f.File.Close()
	if err != nil {
		return err
	}

	newF, err := f.session.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE)
	if err != nil {
		return err
	}

	f.File = newF
	return nil
}

// Remove a file that is open via SFTP.
func (f *File) Remove() error {
	f.File.Close()
	return f.session.Remove(f.Name())
}
