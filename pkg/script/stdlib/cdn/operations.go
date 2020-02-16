package cdn

import (
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
)

func (env Environment) OpenFile(path string) file.Type {
	return file.New(&File{
		Path:       path,
		Uploader:   env.Uploader,
		Downloader: env.Downloader,
	})
}

func (env Environment) openFile(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	return env.OpenFile(path), nil
}
