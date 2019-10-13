package assets

import (
	"net/http"

	"github.com/kcarretto/paragon/script"
)

// Import the assets library to enable scripts to load assets from the provided filesystem.
func Import(assets http.FileSystem) script.Library {
	return script.Library{
		"load": Load(assets),
	}
}
