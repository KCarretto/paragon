package agent

import (
	"bytes"
	"context"
	"io"
	"time"

	"go.uber.org/zap"
)

// drain all available output messages into a buffer
func drain(ctx context.Context, output <-chan []byte) io.Reader {
	// TODO: Configure buffer size
	buffer := bytes.NewBuffer(make([]byte, 0, len(output)*100))
	for {
		select {
		case <-ctx.Done():
			return buffer
		case msg := <-output:
			_, err := buffer.Write(msg)
			// TODO: Handle error
			if err != nil {
				panic(err)
			}
		default:
			return buffer
		}
	}
}

func (agent Agent) writer(ctx context.Context) {
	lastWritten := time.Unix(0, 0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := drain(ctx, agent.output)

			for _, meta := range agent.transports.List() {
				// TODO: Defer recover in case write panics

				// TODO: Configure transport logger
				transport, err := agent.transports.Get(meta.Name, zap.NewNop(), agent)
				if err != nil {
					// TODO: Log error
					continue
				}
				delay := meta.Interval - time.Since(lastWritten)
				time.Sleep(delay)

				var buf []byte
				n, err := io.CopyBuffer(transport, data, buf)
				if err != nil {
					// TODO: Log error
					if err := agent.transports.CloseTransport(meta.Name); err != nil {
						// TODO: Log error
					}
					continue
				}
				if n > 0 {
					lastWritten = time.Now()
				}
			}

			// TODO: ErrNoTransportAvailable
		}
	}
}
