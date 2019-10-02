package transport

import (
	"io"
	"sync"
)

// ResultWriter buffers result output that will be transported to the server.
type ResultWriter interface {
	io.Writer
	WriteResult(*Result)
}

// A PayloadWriter sends payloads from the server to a consumer.
type PayloadWriter interface {
	WritePayload([]byte)
	// WriteError(error)
}

// A Receiver consumes payloads from the server.
type Receiver struct {
	mu sync.Mutex

	Decoder
	Handler func(Payload, error)
}

// WritePayload implements PayloadWriter, passing data to the handler for consumption. It is safe
// for concurrent use by multiple transports.
func (recv *Receiver) WritePayload(data []byte) {
	recv.mu.Lock()
	defer recv.mu.Unlock()

	recv.Handler(recv.Decode(data))
}
