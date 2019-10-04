package transport

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

// A TimestampWriterTo extends io.WriterTo with a method that returns the timestamp of when the
// contents were last successfully written.
type TimestampWriterTo interface {
	io.WriterTo
	Timestamp() time.Time
}

// Buffer is a write safe io.Writer & io.WriterTo that is used to buffer output and safely copy it
// to a transport writer. If copying fails, buffer does not lose data.
type Buffer struct {
	Encoder

	Metadata    Metadata
	MaxIdleTime time.Duration

	mu        sync.RWMutex
	output    *bytes.Buffer
	results   []Result
	timestamp time.Time
}

// WriteResult writes structured task execution output to the response.
func (b *Buffer) WriteResult(result Result) {
	b.mu.Lock()
	b.results = append(b.results, result)
	b.mu.Unlock()
}

// Write buffers the provided data until it is consumed by a transport. It is safe for
// concurrent use.
func (b *Buffer) Write(data []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.output == nil {
		b.output = bytes.NewBuffer(data)
		return len(data), nil
	}

	return b.output.Write(data)
}

// WriteTo implements io.WriterTo, allowing Buffer to be used by io.Copy. It will write it's
// contents to the provided writer. It will preserve any data in the case of error.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	output := b.copyOutput()
	results := b.copyResults()

	if len(output) <= 0 && len(results) <= 0 && time.Since(b.Timestamp()) < b.MaxIdleTime {
		return 0, nil
	}

	resp := Response{
		Metadata: b.Metadata,
		Log:      output,
		Results:  results,
	}
	data, err := b.Encode(resp)
	if err != nil {
		b.restore(results, output)
		return 0, fmt.Errorf("failed to encode response data: %w", err)
	}

	n, err := w.Write(data)
	if err != nil {
		b.restore(results, output)
		return 0, fmt.Errorf("failed to transport response data: %w", err)
	}

	// Update timestamp
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

// copyResults returns a copy of the current result buffer and truncates it.
func (b *Buffer) copyResults() (results []Result) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.results) <= 0 {
		return
	}

	results = make([]Result, len(b.results))
	copy(results, b.results)
	b.results = b.results[len(b.results):]

	return
}

// copyOutput returns a copy of the current output buffer and truncates it.
func (b *Buffer) copyOutput() (data []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.output == nil {
		b.output = &bytes.Buffer{}
	}

	if b.output.Len() <= 0 {
		return
	}

	data = make([]byte, b.output.Len())
	copy(data, b.output.Bytes())
	b.output.Reset()

	return
}

// restore the provided state to the buffer.
func (b *Buffer) restore(results []Result, data []byte) error {
	// Restore results buffer
	for _, result := range results {
		b.WriteResult(result)
	}

	// Restore output buffer
	size := len(data)
	n, err := b.Write(data)
	if err == nil && n != size {
		err = fmt.Errorf("failed to restore full output buffer")
	}

	return err
}

// NewBuffer Initializes and returns a new transport buffer.
func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Encoder: NewDefaultEncoder(),
		output:  bytes.NewBuffer(data),
	}
}
