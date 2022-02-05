package main

const dropperMain = `package main

import (
	"bytes"
	"context"

	assetslib "github.com/kcarretto/paragon/pkg/script/stdlib/assets"
	"github.com/kcarretto/paragon/pkg/drop"
)

func main() {
	ctx := context.Background()

	bundle := assetslib.TarGZBundler{
		Buffer: bytes.NewBuffer(Assets),
	}

	fs, err := bundle.FileSystem()
	if err != nil {
		panic(err)
	}

	drop.TheBase(ctx, fs)
}
`
