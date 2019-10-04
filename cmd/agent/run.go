package main


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
	handler := func(w transport.ResultWriter, payload transport.Payload, err error) {
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
	}

	// Initialize Receiver
	receiver := &transport.Receiver{
		ResultWriter: buffer,
		Decoder:      transport.NewDefaultDecoder(),
		Handler:      transport.PayloadHandlerFn(handler),
	}

	// Register Transports
	addTransports(receiver, registry)

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
	}
}
