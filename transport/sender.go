package transport

import (
	"io"
	"time"
)

type MultiSender struct {
	Transports *Registry
	OnError    func(Transport, error)
}

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
