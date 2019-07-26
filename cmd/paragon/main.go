package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kcarretto/paragon/internal/engine"
)

func main() {
	fmt.Println("Starting")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	agent := engine.New()
	agent.Run(ctx)

	fmt.Println("Exiting")
}

// )

// // TODO: Define this type
// type result string

// // TODO: Define this type
// type task interface {
// 	Execute(ctx context.Context) result
// }

// type agent struct {
// 	results chan result
// 	done    sync.WaitGroup
// }

// // TODO: Determine an actual way to log things.
// func log(msg string) {
// 	fmt.Println(msg)
// }

// // Run the agent until the context is cancelled or it's timeout is reached.
// func (a *agent) Run(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log(ctx.Err().Error())
// 			a.done.Wait()
// 			return
// 		default:
// 			// TODO: Enable configuration of drain timeout
// 			drainCtx, cancel := context.WithTimeout(ctx, time.Second*5)
// 			defer cancel()
// 			results := a.drainResults(drainCtx)
// 			tasks := a.report(results)
// 			a.execute(tasks)
// 			log("Sleeping 3s")
// 			time.Sleep(time.Second * 3)
// 		}
// 	}
// }
// func (a *agent) drainResults(ctx context.Context) []result {
// 	var results []result
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return results
// 		case r := <-a.results:
// 			results = append(results, r)
// 		default:
// 			return results
// 		}
// 	}
// }

// func (a *agent) execute(tasks []task) {
// 	log("Executing tasks")
// 	for _, t := range tasks {
// 		a.done.Add(1)
// 		go func(script task) {
// 			r := script.Execute()
// 			a.results <- r
// 			a.done.Done()
// 		}(t)
// 	}
// }

// func (a *agent) report(results []result) []task {
// 	log("Reporting results")
// 	for i, r := range results {
// 		fmt.Printf("Result %d: %+v\n", i, r)
// 	}
// 	return []task{}
// }

// func main() {
// 	log("Starting")
// 	paragon := &agent{}
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
// 	defer cancel()
// 	paragon.Run(ctx)
// 	log("Exiting")
// }
