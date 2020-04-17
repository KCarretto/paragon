package crypto

import "github.com/kcarretto/paragon/pkg/script"

// Library prepares a new crypto library for use within a script environment.
func Library() script.Library {
	return script.Library{
		"generateKey": script.Func(generateKey),
		"encrypt":     script.Func(encrypt),
		"decrypt":     script.Func(decrypt),
	}
}

// Include the crypto library in a script environment.
func Include() script.Option {
	return script.WithLibrary("crypto", Library())
}
