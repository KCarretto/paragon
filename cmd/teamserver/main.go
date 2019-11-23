package main

import (
	"context"

	"github.com/kcarretto/paragon/pkg/teamserver"
)

func main() {
	ctx := context.Background()

	client := getClient(ctx)
	defer client.Close()

	server := teamserver.Server{
		Log:       newLogger(),
		EntClient: client,
	}
	server.Run()
}
