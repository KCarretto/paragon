package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/script/stdlib"
	"github.com/kcarretto/paragon/transport"
	"go.uber.org/zap"
)

func run(ctx context.Context, logger *zap.Logger) {
	// Initialize buffer
	buffer := &transport.Buffer{
		Encoder:     transport.NewDefaultEncoder(),
		Metadata:    transport.Metadata{},
		MaxIdleTime: time.Second * 5,
	}

	// Configure the logger
	configureLogger(logger, buffer)

	// Initialize registry
	registry := &transport.Registry{}

	// Handle payloads from the server
	hLogger := logger.Named("scripts")
	handler := func(w transport.ResultWriter, payload transport.Payload, err error) {
		if err != nil {
			hLogger.Error("Payload decode error", zap.Error(err))
			return
		}

		hLogger.Debug("Executing new payload from server",
			zap.Int("num_tasks", len(payload.Tasks)),
			zap.Reflect("payload", payload),
		)

		for _, task := range payload.Tasks {
			output := transport.NewResult(task)
			code := script.New(
				task.ID,
				bytes.NewBuffer(task.Content),
				script.WithOutput(output),
				stdlib.Load(),
			) // TODO: Add libraries, set output
			fmt.Println(code.Libraries)
			// TODO: Context timeout
			if err := code.Exec(ctx); err != nil {
				hLogger.Error("failed to execute script", zap.Error(err), zap.String("task_id", task.ID))
			} else {
				hLogger.Debug("completed script execution", zap.String("task_id", task.ID))
			}

			output.CloseWithError(err)
			w.WriteResult(output)
		}
	}

	// Initialize Receiver
	receiver := &transport.Receiver{
		ResultWriter: buffer,
		Decoder:      transport.NewDefaultDecoder(),
		Handler:      transport.PayloadHandlerFn(handler),
	}

	// Register Transports
	addTransports(logger.Named("transport"), receiver, registry)

	// Initialize MultiSender
	sender := transport.MultiSender{
		Transports: registry,
		OnError: func(t transport.Transport, err error) {
			logger.Named("transport").Named(t.Name).Error("Failed to transport data", zap.Error(err))
		},
	}

	// Flush buffer using multi sender
	for {
		if _, err := sender.Send(buffer); err != nil {
			// TODO: Handle encode error
			// TODO: Handle ErrNoTransports
		}
		time.Sleep(time.Millisecond * 50)
	}
}
