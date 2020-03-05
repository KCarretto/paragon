package worker

import (
	"bytes"
	"context"
	"log"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	cdnlib "github.com/kcarretto/paragon/pkg/script/stdlib/cdn"
	filelib "github.com/kcarretto/paragon/pkg/script/stdlib/file"
	sshlib "github.com/kcarretto/paragon/pkg/script/stdlib/ssh"
)

const ServiceTag = "svc-pg-worker"

type Worker struct {
	cdn.Uploader
	cdn.Downloader
	Graph graphql.Client
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

	code := script.New(
		string(task.ID),
		bytes.NewBufferString(task.Content),
		script.WithOutput(output),
		sshEnv.Include(),
		cdnEnv.Include(),
		filelib.Include(),
		assert.Include(),
	)
	execErr = code.Exec(ctx)
}
