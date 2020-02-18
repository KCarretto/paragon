package assets

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
)

// OpenFile that was packed into the compiled binary. The resulting file does not support many
// operations such as Chown, Write, etc. but you may read it's contents or copy it to another file
// i.e. one opened by ssh or sys.
//
//go:generate go run ../gendoc.go -lib assets -func openFile -param path@String -retval file@File -retval err@Error -doc "OpenFile that was packed into the compiled binary. The resulting file does not support many operations such as Chown, Write, etc. but you may read it's contents or copy it to another file i.e. one opened by ssh or sys."
//
// @callable:	assets.openFile
// @param:		path 	@String
// @retval:		file 	@File
// @retval:		err 	@Error
//
// @usage:		f, err = assets.openFile("/path/to/asset")
func (env *Environment) OpenFile(path string) (file.Type, error) {
	if env == nil || env.Assets == nil {
		return file.Type{}, fmt.Errorf("no assets available")
	}
	f, err := env.Assets.Open(path)
	if err != nil {
		return file.Type{}, err
	}

	return file.New(File{f, path}), nil
}

func (env *Environment) openFile(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	retVal, retErr := (*Environment).OpenFile(env, path)
	return script.WithError(retVal, retErr), nil
}
