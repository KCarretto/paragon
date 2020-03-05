package cdn

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/script"
)

// Environment used to configure the behaviour of calls to the cdn library.
type Environment struct {
	cdn.Uploader
	cdn.Downloader
}

// Library prepares a new cdn library for use within a script environment.
func (env *Environment) Library(options ...func(*Environment)) script.Library {
	if env == nil {
		panic(fmt.Errorf("cannot include cdn library without setting non-nil environment"))
	}

	return script.Library{
		"openFile": script.Func(env.openFile),
	}
}

// Include the cdn library in a script environment.
func (env *Environment) Include(options ...func(*Environment)) script.Option {
	return script.WithLibrary("cdn", (*Environment).Library(env, options...))
}
