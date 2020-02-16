package worker

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/script"
	cdnlib "github.com/kcarretto/paragon/pkg/script/stdlib/cdn"
	sshlib "github.com/kcarretto/paragon/pkg/script/stdlib/ssh"

	// "github.com/kcarretto/paragon/pkg/script/workerlib"

	"golang.org/x/crypto/ssh"
)

const ServiceTag = "svc-worker"

type credStore map[int]map[int]*ent.Credential

func (store credStore) addCredentialForTarget(targetID int, credential *ent.Credential) {
	if credential == nil {
		return
	}

	targetStore, ok := store[targetID]
	if !ok || targetStore == nil {
		targetStore = make(map[int]*ent.Credential)
		store[targetID] = targetStore
	}

	targetStore[targetID] = credential
}

func (store credStore) AddCredentials(targetID int, credentials ...*ent.Credential) {
	if store == nil {
		log.Printf("[ERR] Cannot add credentials to null store")
		return
	}
	for _, creds := range credentials {
		log.Printf("[DBG] Adding credential to target %d: %+v", targetID, creds)
		store.addCredentialForTarget(targetID, creds)
	}
	log.Printf("[DBG] New credentials store: %+v", store)
}

func (store credStore) ConfigureSSH(targetID int) (configs []*ssh.ClientConfig) {
	if store == nil {
		return
	}

	targetCreds, ok := store[targetID]
	if !ok || targetCreds == nil {
		return
	}

	creds := make(map[string][]ssh.AuthMethod)
	for _, credential := range targetCreds {
		if credential == nil {
			continue
		}

		userCreds, ok := creds[credential.Principal]
		if !ok || userCreds == nil {
			userCreds = []ssh.AuthMethod{}
			creds[credential.Principal] = userCreds
		}

		// TODO: Handle pubkey/privkey credential type
		creds[credential.Principal] = append(userCreds, ssh.Password(credential.Secret))
	}

	for user, authMethods := range creds {
		configs = append(configs,
			&ssh.ClientConfig{
				User:            user,
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Auth:            authMethods,
			})
	}

	return
}

type Worker struct {
	cdn.Uploader
	cdn.Downloader
	SSH   *sshlib.Connector
	Graph graphql.Client

	credStore
}

func (w *Worker) HandleTaskQueued(ctx context.Context, info event.TaskQueued) {
	log.Printf("Handling new task queued event")

	target := info.Target
	if target == nil || target.ID == 0 {
		log.Printf("[DBG] task queued event with invalid target")
		return
	}

	if w.credStore == nil {
		w.credStore = make(credStore)
	}

	w.AddCredentials(target.ID, info.Credentials...)

	task := info.Task
	if task == nil || task.ID == 0 {
		// TODO: Log invalid task
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

	w.ExecTargetTask(ctx, task, target)
}

func (w *Worker) ExecTargetTask(ctx context.Context, task *ent.Task, target *ent.Target) {
	log.Printf("[DBG] Executing new task (%d)", task.ID)
	if w.SSH == nil {
		w.SSH = &sshlib.Connector{}
	}

	configs := w.ConfigureSSH(target.ID)
	for i, config := range configs {
		log.Printf("[DBG] Adding SSH Client Config (%d): %+v", i+1, config)
	}
	w.SSH.SetConfigs(target.PrimaryIP, configs...)

	sshEnv := &sshlib.Environment{
		RemoteHost: target.PrimaryIP,
		Connector:  w.SSH,
	}
	cdnEnv := &cdnlib.Environment{
		Uploader:   w,
		Downloader: w,
	}

	output := new(bytes.Buffer)
	start := time.Now()
	code := script.New(
		string(task.ID),
		bytes.NewBufferString(task.Content),
		script.WithOutput(output),
		sshEnv.Include(),
		cdnEnv.Include(),
		// sshlib.
		// 	workerlib.Load(
		// 	workerlib.WithSSH(sshlib.Environment{
		// 		RemoteHost: target.PrimaryIP,
		// 		Connector:  w.SSH,
		// 		Downloader: w,
		// 		Uploader:   w,
		// 	}),
		// ),
	)

	var errStr string
	if err := code.Exec(ctx); err != nil {
		errStr = err.Error()
	}
	end := time.Now()

	outputStr := output.String()
	if err := w.Graph.SubmitTaskResult(ctx, models.SubmitTaskResultRequest{
		ID:            task.ID,
		Output:        &outputStr,
		Error:         &errStr,
		ExecStartTime: &start,
		ExecStopTime:  &end,
	}); err != nil {
		log.Printf("[ERR] Failed to submit task execution result: %+v", err)
	}
}
