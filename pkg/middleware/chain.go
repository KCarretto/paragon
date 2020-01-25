package middleware

import "net/http"

// Chain is used to chain multiliple middleware calls for http handlers
func Chain(handler http.Handler, middleware ...func(http.Handler) http.HandlerFunc) http.Handler {
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
