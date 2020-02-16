package assets

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
)

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

	return (*Environment).OpenFile(env, path)
}
