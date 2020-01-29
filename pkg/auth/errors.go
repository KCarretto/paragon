package auth

import "fmt"

var (
	ErrNotAuthenticated = fmt.Errorf("valid authentication for actor is required")
	ErrInvalidSignature = fmt.Errorf("valid signature for service is required")
	ErrMissingHeaders   = fmt.Errorf("missing authentication header(s)")
	ErrUnauthorized     = fmt.Errorf("actor is unauthorized to perform requested action")
)
