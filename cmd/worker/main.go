package main

import (
	"fmt"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/script/workerlib/ssh"
	"github.com/kcarretto/paragon/pkg/worker"
)

func main() {
	teamserverURL := "http://127.0.0.1:80"
	if url := os.Getenv("TEAMSERVER_URL"); url != "" {
		teamserverURL = url
	}
	graph := graphql.Client{
		URL: fmt.Sprintf("%s/%s", teamserverURL, "graphql"),
	}

	cdnURL := teamserverURL
	if url := os.Getenv("CDN_URL"); url != "" {
		cdnURL = url
	}
	cdn := cdn.CDN{URL: cdnURL}

	log.Printf("Testing worker")

	w := &worker.Worker{
		Uploader:   cdn,
		Downloader: cdn,
		SSH:        &ssh.Connector{},
		Graph:      graph,
	}

	sub, _ := newSubscriber(context.Background())
	if closer, ok := sub.(io.Closer); ok {
		defer closer.Close()
	}

	topic := os.Getenv("PUB_TOPIC")
	if topic == "" {
		log.Println("[WARN] No PUB_TOPIC environment variable set to publish events, is this a mistake?")
	}

	taskHandler := func(ctx context.Context, data []byte) {
		log.Printf("[INFO] message recieved: %v\n", string(data))

		var taskQueuedEvent event.TaskQueued
		if err := json.Unmarshal(data, &taskQueuedEvent); err != nil {
			log.Printf("[ERR] Failed to parse event json: %v", err)
			return
		}

		w.HandleTaskQueued(ctx, taskQueuedEvent)
	}

	if err := sub.Subscribe(topic, taskHandler); err != nil {
		panic(err)
	}

	// w.HandleTaskQueued(context.Background(), event.TaskQueued{
	// 	Target: &ent.Target{
	// 		ID:        1,
	// 		PrimaryIP: "127.0.0.1:22",
	// 	},
	// 	Credentials: []*ent.Credential{
	// 		&ent.Credential{
	// 			ID:        11,
	// 			Principal: "root",
	// 			Secret:    "changeme",
	// 		},
	// 	},
	// 	Task: &ent.Task{
	// 		ID:      21,
	// 		Content: simpleTask,
	// 	},
	// 	Tags: []*ent.Tag{
	// 		&ent.Tag{
	// 			ID:   22,
	// 			Name: worker.ServiceTag,
	// 		},
	// 	},
	// })

	for {
		time.Sleep(time.Second)
	}
}
