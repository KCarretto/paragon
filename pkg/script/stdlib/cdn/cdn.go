package cdn

import (
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/script"
)

// Environment used to configure the behaviour of calls to the cdn library.
type Environment struct {
	cdn.Uploader
	cdn.Downloader
}

// Library prepares a new cdn library for use within a script environment.
func Library(uploader cdn.Uploader, downloader cdn.Downloader) script.Library {
	env := &Environment{
		Uploader:   uploader,
		Downloader: downloader,
	}

	return script.Library{
		"openFile": script.Func(env.openFile),
	}
}

// Include the cdn library in a script environment.
func Include(uploader cdn.Uploader, downloader cdn.Downloader) script.Option {
	return script.WithLibrary("cdn", Library(uploader, downloader))
}
