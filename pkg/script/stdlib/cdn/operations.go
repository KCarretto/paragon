package cdn

import (
	"fmt"
	"io"
	"os"

	"github.com/kcarretto/paragon/pkg/script"
	libfile "github.com/kcarretto/paragon/pkg/script/stdlib/file"

	"github.com/spf13/afero"
)

//go:generate go run ../gendoc.go -lib cdn -func upload -param f@File -retval err@Error -doc "Upload a file to the CDN, overwriting any previously stored contents."
func (env Environment) upload(parser script.ArgParser) (script.Retval, error) {
	fd, err := libfile.ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}

	f, err := fd.Open(fd.Path)
	if err != nil {
		return script.WithError(nil, fmt.Errorf("failed to open source file: %w", err)), nil
	}
	defer f.Close()

	if err := env.Upload(fd.Path, f); err != nil {
		return script.WithError(nil, fmt.Errorf("failed to upload: %w", err)), nil
	}

	return nil, nil
}

//go:generate go run ../gendoc.go -lib cdn -func download -param name@String -retval f@File -retval err@Error -doc "Download a file from the CDN."
func (env Environment) download(parser script.ArgParser) (script.Retval, error) {
	name, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	content, err := env.Download(name)
	if err != nil {
		return script.WithError(nil, fmt.Errorf("failed to download file from CDN: %w", err)), nil
	}

	// TODO: Put this fs into the env for cache usage
	fs := afero.NewMemMapFs()
	f, err := fs.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return script.WithError(nil, fmt.Errorf("failed to allocate in-memory buffer for file: %w", err)), nil
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return script.WithError(nil, fmt.Errorf("failed to write file to in-memory buffer: %w", err)), nil
	}

	return script.WithError(
		&libfile.File{
			Path: name,
			Fs:   fs,
		},
		nil,
	), nil
}
