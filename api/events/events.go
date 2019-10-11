// Package events defines a standardized format for paragon events that are published and may be
// subscribed to.
package events

// VERSION describes the currently built version of events.
const VERSION = "0.0.1"

// TODO: Add python support via --python_out=.

//go:generate protoc -I=../../ -I=. --go_out=paths=source_relative:. task.proto
