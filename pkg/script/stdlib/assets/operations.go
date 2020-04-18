package assets

import (
	"github.com/kcarretto/paragon/pkg/script"
	libfile "github.com/kcarretto/paragon/pkg/script/stdlib/file"
)

//go:generate go run ../gendoc.go -lib assets -func file -param path@String -retval f@File -doc "Prepare a descriptor for a file that was packaged into the binary. The descriptor may be used with the file library."
func (env *Environment) file(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	return &libfile.File{
		Path: path,
		Fs:   env.Assets,
	}, nil
}

// Require will be used in the init function for the worker to specify which files you wish to
// include in the asset bundle which will be accessible on the target. Will fatal if error occurs.
//
//go:generate go run ../gendoc.go -lib assets -func require -param filePath@String -doc "Require will be used in the init function for the worker to specify which files you wish to include in the asset bundle which will be accessible on the target. Will fatal if error occurs."
//
// @callable:	assets.require
// @param:		filePath 	@String
//
// @usage:		assets.require("file_on_cdn")
func (env *Environment) Require(filePath string) (err error) {
	f, err := env.Downloader.Download(filePath)
	if err != nil {
		return err
	}
	env.Files = append(env.Files, NamedReader{Reader: f, Name: filePath})
	return nil

}

func (env *Environment) require(parser script.ArgParser) (script.Retval, error) {
	filePath, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	return nil, env.Require(filePath)
}
