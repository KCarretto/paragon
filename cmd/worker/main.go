package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const simpleTask = `
load("ssh", "exec", "sendToTarget", "recvFromTarget")

def main():
	print("Running!")
	exec("touch /tmp/hello_world.txt")

	print(exec("echo hello world"))

	sendToTarget("some_file", "/tmp/some_file")

	recvFromTarget("hello_world", "/etc/ssh/sshd_config")
`

type testDownloader struct {
	File string
}

func (t *testDownloader) Download(filename string) (io.Reader, error) {
	return bytes.NewBufferString(t.File), nil
}

type testUploader struct {
	Result []byte
}

func (t *testUploader) Upload(name string, file io.Reader) (err error) {
	t.Result, err = ioutil.ReadAll(file)
	return
}

func main() {

	sub, _ := newSubscriber(context.Background())
	if closer, ok := sub.(io.Closer); ok {
		defer closer.Close()
	}

	topic := os.Getenv("PUB_TOPIC")
	if topic == "" {
		log.Println("[WARN] No PUB_TOPIC environment variable set to publish events, is this a mistake?")
	}

	if err := sub.Subscribe(topic, func(ctx context.Context, b []byte) {
		log.Printf("[INFO] message recieved: %v\n", string(b))
	}); err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second)
	}

	// log.Printf("Testing worker")
	// uploader := &testUploader{}

	// w := &worker.Worker{
	// 	Uploader: uploader,
	// 	Downloader: &testDownloader{
	// 		File: "test file\n\twoot \twoot!\n\n",
	// 	},
	// 	SSH: &ssh.Connector{},
	// }

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

	// log.Printf("[DBG] File content received: %s", string(uploader.Result))
}
