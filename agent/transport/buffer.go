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

// Timestamp returns the time that the buffer was last successfully written to a transport.
func (b *Buffer) Timestamp() time.Time {
	b.mu.RLock()
	t := b.timestamp
	b.mu.RUnlock()
	return t
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
	b.mu.Lock()

	// Don't waste allocations if there's nothing to write
	size := b.buffer.Len()
	if size <= 0 {
		b.mu.Unlock()
		return 0, io.EOF
	}

	// Default no-op logger if none has been configured
	if b.Logger == nil {
		b.Logger = zap.NewNop()
	}

	// Copy the buffer and truncate it
	data := make([]byte, size)
	copy(data, b.buffer.Bytes())
	b.buffer.Reset()

	// It's important that the lock is released when control is handed to the writer, in case it
	// logs more messages into the buffer.
	b.mu.Unlock()

	// Write the data to the transport
	n, err := w.Write(data)

	// Restore buffer on error
	if err != nil {
		b.mu.Lock()
		num, restoreErr := b.buffer.Write(data)
		b.mu.Unlock()

		// Should never fail to restore buffer
		if restoreErr != nil {
			// Note: Cannot call logger inside lock, as it may write to buffer
			b.Logger.DPanic(
				"Failed restore buffer",
				zap.NamedError("transport_error", err),
				zap.Error(restoreErr),
			)
			return 0, err
		}

		// Should always restore full size of buffer, otherwise this leads to log loss.
		if num != size {
			// Note: Cannot call logger inside lock, as it may write to buffer
			b.Logger.DPanic(
				"Failed to restore entire buffer",
				zap.Int("written_bytes", num),
				zap.Int("total_bytes", size),
			)
		}

		return 0, err
	}

	// Need to aquire lock one last time to update timestamp if bytes were written
	if n > 0 {
		b.mu.Lock()
		b.timestamp = time.Now()
		b.mu.Unlock()
	}

	return int64(n), nil
}

// Sync is a nop to implement zapcore.WriteSyncer
func (b *Buffer) Sync() error {
	return nil
}
