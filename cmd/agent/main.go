package main

import (
	"context"
	"go.uber.org/zap"
	"github.com/kcarretto/paragon/transport"
)

func getLogger(*transport.Buffer) *zap.Logger {
	return zap.NewNop()
}

func main() {
	// Initialize context
	ctx := context.Background()

	// Initialize buffer
	buffer := &transport.Buffer{
		MaxIdleTime: time.Minute * 1,
	}

	// Initialize logger
	logger := getLogger(buffer)

	// Initialize registry
	registry := &transport.Registry{}

	// Initialize Interpreter
	interpreter := script.NewInterpreter()
	// stdlib.Compile(
	// stdlib.WithRegistry(registry),

	// Initialize Receiver
	receiver := &transport.Receiver{
		Decoder: transport.NewDefaultDecoder(),
		Handler: func(payload transport.Payload, err error) {
			for _, task := range payload.Tasks {
				output := transport.NewResult(task)
				// TODO: Context timeout
				if err := interpreter.Execute(ctx, task.Content, output); err != nil {
					// TODO: Handle task execution error
				}
				output.CloseWithError(err)
				buffer.WriteResult(output)
			}
		},
	}

	
	// Register Transports
	registry.Add(transport.New(
		"http",
		&http.Transport{
			// Inject transport dependencies
			PayloadWriter: receiver,
			Logger: logger.Named("transport").Named("http"),
		},
	))
	
	// Initialize MultiSender
	sender := transport.NewMultiSender(registry, func(meta transport.Meta, err error) {
			// HANDLE error
				* Transport error
	})
	
	// Flush buffer using multi sender
	for {
			select {
				// TODO: Handle SIGINT
				// TODO: Handle SIGTERM
				default:
				if err := buffer.Flush(sender); err != nil {
					// HANDLE errors
						* Encode error
						* ErrNoTransports
				}
			}
	}
	
}