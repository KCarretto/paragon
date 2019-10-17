package agent

import "fmt"

// ErrNoTransports occurs if all transports fail to send a message or if none are available.
var ErrNoTransports = fmt.Errorf("all available transports failed to send message")
