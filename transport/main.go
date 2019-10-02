package transport

// package agent

// import (
// 	"io"
// )

// func main() {
// 	// // Initialize context
// 	// ctx := context.Background()
// 	//
// 	// // Initialize buffer
// 	// buffer := transport.NewBuffer(
// 	// 		transport.MaxIdle(time.Minute * 1),
// 	// )
// 	//
// 	// // Initialize logger
// 	// logger := getLogger(buffer)
// 	//
// 	// // Initialize registry
// 	// registry := transport.NewRegistry()
// 	//
// 	// // Initialize Interpreter
// 	// interpreter := script.NewInterpreter(logger.Named("executor"), stdlib.Compile(
// 	//		stdlib.WithRegistry(registry),
// 	// ))
// 	//
// 	// // Initialize Receiver
// 	// receiver := agent.NewReceiver(buffer, func(writer agent.ResponseWriter, payload *agent.Payload) {
// 	//		for _, task := range payload.Tasks {
// 	//			output := writer.NewResult(task)
// 	//			// TODO: Context timeout
// 	//			if err := interpreter.Execute(ctx, task.Content, output); err != nil {
// 	//				// HANDLE error
// 	//					* Task execution error
// 	//			}
// 	//			output.Complete(err)
// 	//		}
// 	// })
// 	//
// 	//	// Register Transports
// 	// registry.Add("http", http.Transport{
// 	//		// Inject transport dependencies
// 	//		PayloadWriter: receiver,
// 	//		Logger: logger.Named("transport").Named("http"),
// 	// })
// 	//
// 	// // Initialize MultiSender
// 	// sender := transport.NewMultiSender(registry, func(meta transport.Meta, err error) {
// 	// 		// HANDLE error
// 	//			* Transport error
// 	// })
// 	//
// 	// // Flush buffer using multi sender
// 	// for {
// 	//		select {
// 	//			// TODO: Handle SIGINT
// 	//			// TODO: Handle SIGTERM
// 	//			default:
// 	// 			if err := buffer.Flush(sender); err != nil {
// 	//				// HANDLE errors
// 	//					* Encode error
// 	//					* ErrNoTransports
// 	//			}
// 	//		}
// 	// }
// 	//
// }

// type PayloadWriter interface {
// 	WritePayload(data []byte)
// }

// type ResponseWriter interface {
// 	io.Writer
// 	WriteOutput(taskID string, output []byte)
// }

// type Sender interface {
// 	Send(PayloadWriter, []byte) error
// }

// type Receiver interface {
// 	Receive(ctx context.Context, ResponseWriter, Payload)
// }

// type Decoder interface {
// 	Decode([]byte) (Payload, error)
// }

// type Payload struct {
// 	data []byte
// }

// func (payload *Payload) Write(p []byte) (int, error) {
// 	payload.WritePayload(p)
// 	return len(p), nil
// }

// func (payload *Payload) WritePayload(p []byte) {
// 	payload.data = append(payload.data, p...)
// }

// func (payload *Payload) Bytes() []byte {
// 	return payload.data
// }

// type Agent struct {
// 	ctx context.Context
// 	Decoder
// 	Encoder
// 	Receiver
// 	Sender

// 	Buffer
// }

// func (agent *Agent) WritePayload(data []byte) {
// 	payload, err := agent.Decoder.Decode(data)
// 	if err != nil {
// 		// TODO: Enable error handling
// 		panic(err)
// 	}

// 	go agent.Receiver.Receive(agent.ctx, agent.Buffer, payload)
// }

// func (agent *Agent) Flush(Response) (*Payload, error) {
// 	// Encode response
// 	// Iterate transport registry
// 	// 		sleep interval if required
// 	// 		call buffer.Flush(sender)
// 	//		if no error, return nil
// 	// return ErrNoTransports
// }
