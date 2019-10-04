package local

import (
	"io"
	"os"
	"sync"

	"github.com/kcarretto/paragon/transport"
)

// Transport is a local transport that reads from stdin and writes to stdout.
type Transport struct {
	io.Reader
	transport.PayloadWriter

	wg     *sync.WaitGroup
	input  chan []byte
	closed bool
}

// Write sends response data to stdout, and ensure that stdin is being consumed.
func (t Transport) Write(data []byte) (int, error) {
	t.ensureActive()
	return os.Stdout.Write(data)
}

// Close stops the goroutine responsible for consuming stdin if it is active.
func (t Transport) Close() error {
	if t.input != nil {
		close(t.input)
	}
	t.wg.Wait()

	t.closed = true

	return nil
}

// ensureActive starts a goroutine to consume stdin if the transport is not yet active.
func (t Transport) ensureActive() {
	if !t.closed {
		return
	}

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for msg := range t.input {
			t.WritePayload(msg)
		}
	}()

	t.closed = false
}

// New initializes and returns a local transport, which must be closed.
func New(receiver transport.PayloadWriter) Transport {
	var wg sync.WaitGroup

	input := make(chan []byte, 1)
	t := Transport{
		Reader:        os.Stdin,
		PayloadWriter: receiver,
		wg:            &wg,
		input:         input,
	}

	return t
}
