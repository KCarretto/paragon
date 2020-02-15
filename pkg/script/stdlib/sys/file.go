package sys

import (
	"fmt"
	"os"
	"path"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
)

// File is our impl of File.
type File struct {
	*os.File
}

// Move uses os.Rename to move a file/folder
func (f File) Move(dstPath string) error {
	dir := path.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory %q: %w", dir, err)
	}

	err := os.Rename(f.Name(), dstPath)
	if err != nil {
		return err
	}

	err = f.File.Close()
	if err != nil {
		return err
	}

	newF, err := os.Open(dstPath)
	if err != nil {
		return err
	}

	f.File = newF
	return nil
}

// Remove uses os.Remove to remove a file/folder. WARNING: basically works like rm -rf
func (f File) Remove() error {
	f.Close()
	return os.RemoveAll(f.Name())
}

// OpenFile uses os.Open to Open a file.
func OpenFile(filePath string) (file.Type, error) {
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return file.New(nil), fmt.Errorf("failed to create parent directory %q: %w", dir, err)
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return file.New(nil), err
	}
	return file.New(File{
		f,
	}), nil
}

func openFile(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	return OpenFile(path)
}
