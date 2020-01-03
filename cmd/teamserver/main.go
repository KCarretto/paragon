package main

import (
	"context"
	"github.com/kcarretto/paragon/pkg/teamserver"
	"log"
)

func main() {
	ctx := context.Background()

	client := getClient(ctx)
	defer client.Close()

	server := teamserver.Server{
		Log:       newLogger(),
		EntClient: client,
	}
	log.Println("Starting teamserver")
	server.Run()
}
