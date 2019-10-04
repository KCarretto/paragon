package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/kcarretto/paragon/transport"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		var agentResponse transport.Response
		if err := json.Unmarshal(data, &agentResponse); err != nil {
			panic(err)
		}

		fmt.Printf("\n\nReceived request: %+v\n", agentResponse)

		var tasks []transport.Task
		if rand.Intn(2) == 0 {
			task := transport.Task{
				ID:      "task1",
				Content: []byte(`print("Task from C2 ;)")`),
			}
			tasks = append(tasks, task)
		}

		resp := transport.Payload{
			Tasks: tasks,
		}

		respData, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		if _, err = w.Write(respData); err != nil {
			panic(err)
		}

		fmt.Println("Responded to agent")

	})

	fmt.Println("Running C2. Serving on 127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		panic(err)
	}
}
