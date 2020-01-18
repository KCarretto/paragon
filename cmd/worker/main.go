package main

import (
	"context"
	"log"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/worker"
)

const simpleTask = `
load("ssh", "exec")

def main():
	print("Running!")
	exec("touch /tmp/hello_world.txt")
`

func main() {
	log.Printf("Testing worker")
	w := &worker.Worker{}

	w.HandleTaskQueued(context.Background(), event.TaskQueued{
		Target: &ent.Target{
			ID:        1,
			PrimaryIP: "127.0.0.1:22",
		},
		Credentials: []*ent.Credential{
			&ent.Credential{
				ID:        11,
				Principal: "root",
				Secret:    "changeme",
			},
		},
		Task: &ent.Task{
			ID:      21,
			Content: simpleTask,
		},
		Tags: []*ent.Tag{
			&ent.Tag{
				ID:   22,
				Name: worker.ServiceTag,
			},
		},
	})

}
