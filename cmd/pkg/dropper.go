package main

const dropperMain = `package main

import (
	"context"

	"github.com/kcarretto/paragon/drop"
)

func main() {
	ctx := context.Background()
	drop.TheBase(ctx, Assets)
}

`
