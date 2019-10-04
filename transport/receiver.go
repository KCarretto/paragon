package transport

import (
	"io"
	"sync"
)

// ResultWriter buffers result output that will be transported to the server.
type ResultWriter interface {
	io.Writer
	WriteResult(Result)
}

// A PayloadWriter sends payloads from the server to a consumer.
type PayloadWriter interface {
	WritePayload([]byte)
	// WriteError(error)
}

// A PayloadHandler is responsible for consuming and acting upon payloads received from the server.
type PayloadHandler interface {
	HandlePayload(ResultWriter, Payload, error)
}

// PayloadHandlerFn is a function that implements PayloadHandler.
type PayloadHandlerFn func(ResultWriter, Payload, error)

// HandlePayload wraps fn to provide an implementation of PayloadHandler.
func (fn PayloadHandlerFn) HandlePayload(w ResultWriter, data Payload, err error) {
	fn(w, data, err)
}

// A Receiver consumes payloads from the server.
type Receiver struct {
	mu sync.Mutex

	ResultWriter
	Decoder
	Handler func(ResultWriter, Payload, error)
}

// WritePayload implements PayloadWriter, passing data to the handler for consumption. It is safe
// for concurrent use by multiple transports.
func (recv *Receiver) WritePayload(data []byte) {
	recv.mu.Lock()
	defer recv.mu.Unlock()

	payload, err := recv.Decode(data)
	recv.Handler(recv.ResultWriter, payload, err)
}
