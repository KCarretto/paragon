package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/agent/transport"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		var payload transport.AgentPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			panic(err)
		}

		fmt.Printf("\n\nReceived request: %s\n", payload.Output)

		resp := transport.ServerPayload{
			Tasks: []transport.TaskPayload{
				transport.TaskPayload{
					ID:      "task1",
					Content: []byte(`print("Task from C2 ;)")`),
				},
			},
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
