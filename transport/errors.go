package transport

import "errors"

// ErrNoTransportAvailable occurs when all available transports fail to report output, or if no transports were configured.
var ErrNoTransportAvailable = errors.New("no transport available")
