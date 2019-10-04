package main

import (
	"bytes"
	"context"
	"time"

	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/transport"
	"github.com/kcarretto/paragon/transport/local"
	"go.uber.org/zap"
)

func getLogger(*transport.Buffer) *zap.Logger {
	return zap.NewNop()
}

func main() {
	// Initialize context
	ctx := context.Background()

	// Initialize buffer
	buffer := &transport.Buffer{
		Encoder:     transport.NewDefaultEncoder(),
		Metadata:    transport.Metadata{},
		MaxIdleTime: time.Second * 5,
	}

	// Initialize logger
	logger := getLogger(buffer)

	// Initialize registry
	registry := &transport.Registry{}

	// Initialize Receiver
	receiver := &transport.Receiver{
		ResultWriter: buffer,
		Decoder:      transport.NewDefaultDecoder(),
		Handler: transport.PayloadHandlerFn(func(w transport.ResultWriter, payload transport.Payload, err error) {
			for _, task := range payload.Tasks {
				output := transport.NewResult(task)
				code := script.New(task.ID, bytes.NewBuffer(task.Content)) // TODO: Add libraries, set output
				// TODO: Context timeout
				if err := code.Exec(ctx); err != nil {
					// TODO: Handle task execution error
				}
				output.CloseWithError(err)
				w.WriteResult(output)
			}
		}),
	}

	// Register Transports
	registry.Add(transport.New(
		"local",
		local.New(receiver),
	))

	// Initialize MultiSender
	sender := transport.MultiSender{
		Transports: registry,
		OnError: func(t transport.Transport, err error) {
			logger.Named("transport").Named(t.Name).Error("Failed to transport data", zap.Error(err))
		},
	}

	// Flush buffer using multi sender
	for {
		select {
		// TODO: Handle SIGINT
		// TODO: Handle SIGTERM
		default:
			if _, err := sender.Send(buffer); err != nil {
				// TODO: Handle encode error
				// TODO: Handle ErrNoTransports
			}
		}
	}

}
