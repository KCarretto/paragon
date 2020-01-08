package c2

import (
	types "github.com/gogo/protobuf/types"
	"time"
)

// Start recording results for task execution.
func (m *TaskResult) Start() {
	start := time.Now()

	m.ExecStartTime = &types.Timestamp{
		Seconds: start.Unix(),
		Nanos:   int32(start.Nanosecond()),
	}
}

// CloseWithError marks the end of task execution, and provides any error the task resulted in.
func (m *TaskResult) CloseWithError(err error) {
	stop := time.Now()
	m.ExecStopTime = &types.Timestamp{
		Seconds: stop.Unix(),
		Nanos:   int32(stop.Nanosecond()),
	}
	if err != nil {
		m.Error = err.Error()
	}

}
