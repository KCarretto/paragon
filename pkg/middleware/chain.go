package middleware

import "net/http"

// Wrapper is a helper type for declraing the chain segments
type Wrapper func(http.Handler) http.HandlerFunc

// Chain is used to chain mutliple middleware calls for http handlers
func Chain(handler http.Handler, middleware ...Wrapper) http.Handler {
	if len(middleware) < 1 {
		return handler
	}

	wrapped := handler

	// loop in reverse to preserve middleware order
	for i := len(middleware) - 1; i >= 0; i-- {
		wrapped = middleware[i](wrapped)
	}

	return wrapped
}
