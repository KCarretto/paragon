package transport

import (
	"encoding/json"
	"time"
)

// An Encoder is responsible for marshaling a response to bytes.
type Encoder interface {
	Encode(Response) ([]byte, error)
}

// EncoderFn is a function that implements Encoder.
type EncoderFn func(Response) ([]byte, error)

// Encode wraps fn to provide an implementation of Encoder.
func (fn EncoderFn) Encode(resp Response) ([]byte, error) {
	return fn(resp)
}

// NewDefaultEncoder returns a JSON Encoder.
func NewDefaultEncoder() Encoder {
	return EncoderFn(func(resp Response) ([]byte, error) {
		return json.Marshal(resp)
	})
}

// A Decoder is responsible for unmarshaling received data into a payload struct.
type Decoder interface {
	Decode([]byte) (Payload, error)
}

// DecoderFn is a function that implements Decoder.
type DecoderFn func([]byte) (Payload, error)

// Decode wraps fn to provide an implementation of Decoder.
func (fn DecoderFn) Decode(data []byte) (Payload, error) {
	return fn(data)
}

// NewDefaultDecoder returns a JSON Decoder.
func NewDefaultDecoder() Decoder {
	return DecoderFn(func(data []byte) (payload Payload, err error) {
		err = json.Unmarshal(data, &payload)
		return
	})
}

// Payload holds structured information received from a server.
type Payload struct {
	Tasks []Task
}

// Response holds structured information that will be transported to a server.
type Response struct {
	Metadata Metadata
	Results  []*Result `json:"results"`
	Log      []byte    `json:"log"`
}

// Metadata holds agent metadata useful for identifying an agent.
type Metadata struct{}

// Task stores instructions to execute and metadata.
type Task struct {
	ID      string `json:"id"`
	Content []byte `json:"content"`
}

// Result stores task execution output and metadata.
type Result struct {
	ID     string `json:"id"`
	Output []byte `json:"output"`
	Error  string `json:"error"`

	ExecStartTime time.Time `json:"exec_start_time"`
	ExecStopTime  time.Time `json:"exec_stop_time"`
}

// Write implements io.Writer by appending output to the current result buffer. It is safe for
// concurrent use and always returns len(p), nil.
func (r *Result) Write(p []byte) (int, error) {
	r.Output = append(r.Output, p...)

	return len(p), nil
}

// Close implements io.Closer by wrapping CloseWithError to indicate no error. Always returns nil.
func (r *Result) Close() error {
	r.CloseWithError(nil)
	return nil
}

// CloseWithError sets the result error and execution stop time.
func (r *Result) CloseWithError(err error) {
	if err != nil {
		r.Error = err.Error()
	}
	r.ExecStopTime = time.Now()
}

// NewResult initializes and returns a new Result to store execution output for the provided task.
func NewResult(task Task) *Result {
	return &Result{
		ID:            task.ID,
		ExecStartTime: time.Now(),
	}
}
