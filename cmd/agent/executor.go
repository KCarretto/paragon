package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/kcarretto/paragon/pkg/agent/transport"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	"github.com/kcarretto/paragon/pkg/script/stdlib/http"
	"github.com/kcarretto/paragon/pkg/script/stdlib/process"
	"github.com/kcarretto/paragon/pkg/script/stdlib/regex"
	"github.com/kcarretto/paragon/pkg/script/stdlib/websockets"
	"github.com/kcarretto/paragon/pkg/script/stdlib/sys"
)

// Executor is responsible for executing tasks from the server.
type Executor struct{
	Metadata	transport.AgentMetadata
}

// ExecuteTask runs a renegade script.
func (exec Executor) ExecuteTask(ctx context.Context, output io.Writer, task *transport.Task) error {

	code := script.New(
		fmt.Sprint(task.GetId()),
		bytes.NewBufferString(task.Content),
		script.WithOutput(output),
		sys.Include(),
		http.Include(),
		file.Include(),
		regex.Include(),
		process.Include(),
		assert.Include(),
		websockets.Include(),
	)

	return code.Exec(ctx)
}
