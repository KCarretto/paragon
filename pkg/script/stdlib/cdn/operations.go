package cdn

import (
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
)

// OpenFile stored on the CDN. Writing to the file will cause an upload to the CDN, overwriting any
// previously stored contents. Reading the file will download it from the CDN. Since operations are
// performed lazily, openFile will never error however reading from or writing to the file may.
//
// @callable:	cdn.openFile
// @param:		name 	@string
// @retval:		file 	@File
//
// @usage:		f = cdn.openFile("file_name")
func (env Environment) OpenFile(name string) file.Type {
	return file.New(&File{
		Path:       name,
		Uploader:   env.Uploader,
		Downloader: env.Downloader,
	})
}

func (env Environment) openFile(parser script.ArgParser) (script.Retval, error) {
	name, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	return env.OpenFile(name), nil
}
