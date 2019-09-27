package transport

import (
	"bytes"
	"io"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Buffer is a write safe io.Writer & io.WriterTo that is used to buffer agent output and safely
// copy it to a transport writer. If copying fails, buffer does not lose data.
type Buffer struct {
	Logger    *zap.Logger
	mu        sync.RWMutex
	buffer    *bytes.Buffer
	timestamp time.Time
}

// NewBuffer Initializes and returns a new transport buffer.
func NewBuffer(p []byte) *Buffer {
	return &Buffer{
		buffer: bytes.NewBuffer(p),
	}
}

// Write buffers the provided data until it is consumed by a transport. It is safe for concurrent use.
func (b *Buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.buffer == nil {
		b.buffer = bytes.NewBuffer(p)
		return len(p), nil
	}

	return b.buffer.Write(p)
}

// WriteTo implements io.WriterTo, allowing Buffer to be used by io.Copy. It will write it's
// contents to the provided writer. It will preserve any data that still needs to be written, even
// in the case of error.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	data := b.readCopy()

	n, err := w.Write(data)
	if err != nil {
		b.restore(data)
		return 0, err
	}

	b.timestamp = time.Now()

	return int64(n), nil
}

// Sync is a nop to implement zapcore.WriteSyncer
func (b *Buffer) Sync() error {
	return nil
}

// Timestamp returns the time that the buffer was last successfully written to a transport.
func (b *Buffer) Timestamp() time.Time {
	b.mu.RLock()
	t := b.timestamp
	b.mu.RUnlock()
	return t
}

func (b *Buffer) readCopy() (data []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.buffer.Len() <= 0 {
		return
	}

	data = make([]byte, b.buffer.Len())
	copy(data, b.buffer.Bytes())
	b.buffer.Reset()

	return
}

// restore the buffer to the provided state.
func (b *Buffer) restore(data []byte) {
	size := len(data)

	n, err := b.Write(data)

	if err != nil || n != size {
		// Default no-op logger if none has been configured
		if b.Logger == nil {
			b.Logger = zap.NewNop()
		}
		b.Logger.DPanic(
			"Failed to restore buffer",
			zap.Error(err),
			zap.Int("written_bytes", n),
			zap.Int("total_bytes", size),
		)
	}
}
