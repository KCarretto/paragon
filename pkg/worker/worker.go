package worker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	cdnlib "github.com/kcarretto/paragon/pkg/script/stdlib/cdn"
	envlib "github.com/kcarretto/paragon/pkg/script/stdlib/env"
	filelib "github.com/kcarretto/paragon/pkg/script/stdlib/file"
	sshlib "github.com/kcarretto/paragon/pkg/script/stdlib/ssh"

	"go.starlark.net/starlark"
)

const ServiceTag = "svc-pg-worker"

type Worker struct {
	cdn.Uploader
	cdn.Downloader
	Graph  graphql.Client
	Config string
}

func (w *Worker) HandleTaskQueued(ctx context.Context, info event.TaskQueued) {
	log.Printf("Handling new task queued event")

	target := info.Target
	if target == nil || target.ID == 0 {
		log.Printf("[DBG] task queued event with invalid target")
		return
	}

	task := info.Task
	if task == nil || task.ID == 0 {
		log.Printf("[DBG] task queued event with invalid task")
		return
	}

	tags := info.Tags
	if tags == nil {
		log.Printf("[DBG] task queued event with invalid tags")
		return
	}

	found := false
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		if tag.Name != ServiceTag {
			continue
		}

		found = true
		break
	}
	if !found {
		log.Printf("[DBG] task queued event without worker tags")
		return
	}

	if target.PrimaryIP == "" {
		log.Printf("[DBG] task queued event with invalid target ip")
		return
	}

	w.ExecTargetTask(ctx, task, target, info.Credentials)
}

func (w *Worker) ExecTargetTask(ctx context.Context, task *ent.Task, target *ent.Target, credentials []*ent.Credential) {
	_, err := w.Graph.ClaimTask(ctx, task.ID)
	if err != nil {
		log.Printf("[ERR] Failed to claim task: %v", err)
	}

	output := &taskOutput{
		ID:    task.ID,
		Ctx:   ctx,
		Graph: w.Graph,
	}
	output.Start()
	var execErr error
	defer func() {
		output.Stop(execErr)
	}()

	log.Printf("[DBG] Executing new task (%d) on %s (%d credentials)",
		task.ID,
		target.PrimaryIP,
		len(credentials),
	)

	env := envlib.Environment{
		PrimaryIP:       target.PrimaryIP,
		OperatingSystem: target.Name, // TODO: Add OperatingSystem field to target
	}

	/* Build Assets Bundle */
	code := script.New(
		string(task.ID),
		strings.NewReader(task.Content),
		script.WithOutput(output),
		filelib.Include(),
		assert.Include(),
		env.Include(),
	)
	if _, err := code.Call("init", starlark.Tuple{}); err != nil {
		execErr = fmt.Errorf("failed to initialize assets: %w", err)
		return
	}

	/* Use Config to execute task */
	sshConnector := &SSHConnector{
		Credentials: credentials,
	}
	defer sshConnector.Close()

	sshEnv := &sshlib.Environment{
		RemoteHost: target.PrimaryIP,
		Connector:  sshConnector,
	}
	defer sshEnv.Close()

	cdnEnv := &cdnlib.Environment{
		Uploader:   w,
		Downloader: w,
	}

	var configScript = strings.NewReader(DefaultConfig)
	if w.Config != "" {
		configScript = strings.NewReader(w.Config)
	}

	config := script.New(
		string(task.ID),
		configScript,
		script.WithOutput(output),
		sshEnv.Include(),
		cdnEnv.Include(),
		filelib.Include(),
		assert.Include(),
	)

	var res starlark.Value = starlark.None
	defer func() {
		if _, err := config.Call("worker_exit", starlark.Tuple{res}); err != nil {
			log.Printf("[ERR] worker_exit failed: %v", err)
		}
	}()

	// TODO: pass assets as argument
	res, execErr = config.Call("worker_run", starlark.Tuple{starlark.String(task.Content), starlark.String("TODO: Assets")})
}
