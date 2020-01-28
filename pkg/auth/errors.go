package auth

import "fmt"

var (
	ErrNotAuthenticated = fmt.Errorf("valid authentication for actor is required")
	ErrUnauthorized     = fmt.Errorf("actor is unauthorized to perform requested action")
)
