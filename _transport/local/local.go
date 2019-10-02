package local

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kcarretto/paragon/agent/transport"
	"go.uber.org/zap"
)

// New starts and returns a local transport that reads tasks from stdin and writes output to stdout.
func New(logger *zap.Logger, tasks transport.Tasker) (io.WriteCloser, error) {
	fmt.Printf("Local task >>")
	tasks.QueueTask(fmt.Sprintf("local_task_%v", time.Now().String()), os.Stdin)

	return Transport{
		logger,
		tasks,
	}, nil
}

// Transport is a local transport for experimental purposes.
type Transport struct {
	logger *zap.Logger
	tasks  transport.Tasker
}

// Write agent output to stdout.
func (t Transport) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

// Close the local transport to stop reading tasks.
func (t Transport) Close() error {
	return nil
}
