// Package codec defines a standardized format for messages sent between the server and agents.
// It does not mandate an encoding mechanism for transporting such data. Any format for serializing
// messages sent between the agent and the server (i.e. protobuf or JSON) is acceptable so long as
// both the agent and the server utilize the same format.
package codec

// VERSION describes the currently built version of codec.
const VERSION = "0.0.1"

// TODO: Add python support via --python_out=.

//go:generate protoc -I=../vendor/ -I=../../ -I=. --go_out=paths=source_relative:. agent.proto server.proto
