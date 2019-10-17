package main

// import (
// 	"context"
// 	"net/http"
// 	"os"

// 	"github.com/kcarretto/paragon/c2"
// 	"go.uber.org/zap"
// )

func main() {
	// ctx := context.Background()

	// logger, err := zap.NewDevelopment()
	// if err != nil {
	// 	panic(err)
	// }

	// httpAddr := os.Getenv("HTTP_ADDR")
	// if httpAddr == "" {
	// 	httpAddr = "127.0.0.1:8080"
	// }

	// resultTopic, err := openTopic(ctx, "tasks.result_received")
	// if err != nil {
	// 	logger.Panic("Failed to open pubsub topic", zap.Error(err))
	// }
	// defer resultTopic.Shutdown(ctx)

	// queueTopic, err := openSubscription(ctx, "tasks.queued")
	// if err != nil {
	// 	logger.Panic("Failed to subscribe to pubsub topic", zap.Error(err))
	// }
	// defer queueTopic.Shutdown(ctx)

	// srv := &c2.Server{
	// 	Log:         logger,
	// 	TaskResults: resultTopic,
	// }

	// go func() {

	// }()

	// logger.Info("Started C2 server", zap.String("http_addr", httpAddr))
	// if err := http.ListenAndServe(httpAddr, srv); err != nil {
	// 	logger.Panic("Failed to serve HTTP", zap.Error(err))
	// }

	// router := http.NewServeMux()
	// router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	data, err := ioutil.ReadAll(req.Body)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	var agentResponse transport.Response
	// 	if err := json.Unmarshal(data, &agentResponse); err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Printf("\n\nReceived request: %+v\n", agentResponse)

	// 	var tasks []transport.Task
	// 	if rand.Intn(2) == 0 {
	// 		task := transport.Task{
	// 			ID:      "task1",
	// 			Content: []byte(`print("Task from C2 ;)")`),
	// 		}
	// 		tasks = append(tasks, task)
	// 	}

	// 	resp := transport.Payload{
	// 		Tasks: tasks,
	// 	}

	// 	respData, err := json.Marshal(resp)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if _, err = w.Write(respData); err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println("Responded to agent")

	// })

}
