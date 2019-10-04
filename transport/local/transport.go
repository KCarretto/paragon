package local

import (
	"io"
	"os"
	"sync"

	"github.com/kcarretto/paragon/transport"
)

type Transport struct {
	io.Reader
	transport.PayloadWriter

	wg     *sync.WaitGroup
	input  chan []byte
	closed bool
}

func (t Transport) Write(data []byte) (int, error) {
	t.ensureActive()
	return os.Stdout.Write(data)
}

func (t Transport) Close() error {
	if t.input != nil {
		close(t.input)
	}
	t.wg.Wait()

	t.closed = true

	return nil
}

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
