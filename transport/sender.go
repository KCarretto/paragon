package transport

import (
	"io"
	"time"
)

// MultiSender is used to send a response buffer to the server using transports from the registry.
type MultiSender struct {
	Transports *Registry
	OnError    func(Transport, error)
}

// Send the response buffer to the server using available transports. If a transport fails to send
// the respond, OnError will be invoked if it is set. Transports will be attempted in the order
// defined by the registry sorting mechanism. If all transports are exhausted, ErrNoTransportAvailable
// will be returned.
func (sender MultiSender) Send(buffer TimestampWriterTo) (int64, error) {
	for _, transport := range sender.Transports.List() {
		delay := transport.Interval - time.Since(buffer.Timestamp())
		time.Sleep(delay)

		n, err := buffer.WriteTo(transport)
		if err != nil && err != io.EOF && sender.OnError != nil {
			sender.OnError(transport, err)
			continue
		}

		return n, nil
	}

	return 0, ErrNoTransportAvailable
}
