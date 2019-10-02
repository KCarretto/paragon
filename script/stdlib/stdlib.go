package stdlib

import (
	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/script/stdlib/sys"
	"go.starlark.net/starlark"
)

func Loader() func(_ *starlark.Thread, name string) (starlark.StringDict, error) {
	libs := map[string]starlark.StringDict{
		"sys": sys.Lib.Compile(),
	}

	return func(_ *starlark.Thread, name string) (starlark.StringDict, error) {
		lib, ok := libs[name]
		if !ok {
			return starlark.StringDict{}, script.ErrMissingLibrary
		}
		return lib, nil
	}
}
