// Package codec defines a standardized format for messages sent between the server and agents.
// It does not mandate an encoding mechanism for transporting such data. Any format for serializing
// messages sent between the agent and the server (i.e. protobuf or JSON) is acceptable so long as
// both the agent and the server utilize the same format.
package codec

import (
	"time"

	types "github.com/gogo/protobuf/types"
)

// Start recording results for task execution.
func (m *Result) Start() {
	start := time.Now()

	m.ExecStartTime = &types.Timestamp{
		Seconds: start.Unix(),
		Nanos:   int32(start.Nanosecond()),
	}
}

// Write appends task execution output to the result.
func (m *Result) Write(p []byte) (int, error) {
	m.Output = append(m.Output, string(p))
	return len(p), nil
}

// CloseWithError marks the end of task execution, and provides any error the task resulted in.
func (m *Result) CloseWithError(err error) {
	stop := time.Now()
	m.ExecStopTime = &types.Timestamp{
		Seconds: stop.Unix(),
		Nanos:   int32(stop.Nanosecond()),
	}
	if err != nil {
		m.Error = err.Error()
	}

}
