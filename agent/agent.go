package agent

import (
	"context"
	"io"
	"sync"
)

/*
A Task is uniquely identifiable job that is queued for execution.

Implementation Details:
	* Write implemented to buffer output of task execution until it is read.
 	* Read implemented to return formatted output of of task execution.
	* Close implemented to perform any cleanup at the end of task execution.
	* ID implemented to return the Task's unique identifier.
	* Execute implemented to run the Task.
*/
type Task interface {
	io.ReadWriteCloser

	ID() []byte
	Execute(*Agent) error
}

// Config stores the agent's configuration values. It may be used by various components (i.e. a
// Transport implementation) to control behaviour (i.e. what server(s) to connect to).
type Config struct{}

// A Parser is responsible for converting raw byte input received from the transport into a Task.
type Parser interface {
	Parse(task []byte) (Task, error)
}

/*
A Transport is used to report the output of task execution, as well as provide any new tasks that
need to be prepared for execution.

Implementation Details:
	* Read implemented to return []byte which can be marshaled into a Task by a Parser.
	* Write implemented to write formatted output to a destination.
	* Close implemented to perform any cleanup when a transport is no longer needed. It should always
  	  be called when stopping a transport.
	* Run implemented to start the transport. It may start go routines that run on an interval if
	  required, so long as they are cleaned up by the close method. Read/Write/Close should not be
	  called on the transport until the Run method is invoked. If it returns an error, the transport
      will be closed.
*/
type Transport interface {
	io.ReadWriteCloser
	Run(context.Context, Config) error
}

/*
An Agent holds the various components responsible for receiving, executing, and outputting the
results of various commands.
*/
type Agent struct {
	mu sync.RWMutex
	wg sync.WaitGroup

	Config
	Parser
	Transport
}
