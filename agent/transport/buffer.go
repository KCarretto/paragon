package transport

import (
	"io"
	"sync"
	"time"
)

// Buffer is a write safe io.Writer & io.WriterTo that is used to buffer agent output and safely
// copy it to a transport writer. If copying fails, buffer does not lose data.
type Buffer struct {
	mu        sync.Mutex
	buf       chan []byte
	timestamp time.Time
}

// NewBuffer Initializes and returns a new transport buffer.
func NewBuffer(maxSize int) *Buffer {
	return &Buffer{
		buf: make(chan []byte, maxSize),
	}
}

// Timestamp returns the time that the buffer was last successfully written to a transport.
func (b *Buffer) Timestamp() time.Time {
	b.mu.Lock()
	t := b.timestamp
	b.mu.Unlock()
	return t
}

// Write buffers the provided data until it is consumed by a transport. It is safe for concurrent use.
// It never errors, unless the make or copy builtins panic.
func (b *Buffer) Write(p []byte) (int, error) {
	// TODO: Select + goroutine + waitgroup
	// TODO: Ensure channel is initialized
	// if b.buf == nil {
	// 	b.buf = make(chan []byte, DefaultMessageBacklog)
	// }
	data := make([]byte, len(p))
	n := copy(data, p)
	b.buf <- data
	return n, nil
}

// WriteTo implements io.WriterTo, allowing Buffer to be used by io.Copy. It will write it's
// contents to the provided writer. It will preserve any data that still needs to be written, even
// in the case of error.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	size, msgs := drain(b.buf)
	if size <= 0 {
		return 0, io.EOF
	}

	msg := marshal(msgs...)
	n, err := w.Write(msg)
	if err != nil {
		b.Write(msg) // TODO: b.buf <-msg or b.Write?
		return 0, err
	}

	return int64(n), nil
}

// Sync is a nop to implement zapcore.WriteSyncer
func (b *Buffer) Sync() error {
	return nil
}

// marshal multiple messages into a single byte array.
func marshal(msgs ...[]byte) (msg []byte) {
	for _, entry := range msgs {
		msg = append(msg, entry...)
	}
	return
}

// drain and return all available output messages from the channel
func drain(output <-chan []byte) (size int, msgs [][]byte) {
	for {
		select {
		case msg := <-output:
			size += len(msg)
			msgs = append(msgs, msg) // TODO: Better way than append?
		default:
			return
		}
	}
}
