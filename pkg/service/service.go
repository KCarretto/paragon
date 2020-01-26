package service

import (
	"context"
	"fmt"
	"net/http"
)

type RequestLogger interface {
	LogRequest(*http.Request, error)
}

type Authenticator interface {
	Authenticate(*http.Request) (context.Context, error)
}

type Authorizer interface {
	Authorize(*http.Request) error
}

type ErrorPresenter interface {
	PresentError(http.ResponseWriter, error)
}

type Handler interface {
	Handle(http.ResponseWriter, *http.Request) error
}

func HTTPHandler(handler http.Handler) HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		handler.ServeHTTP(w, r)
		return nil
	}
}

type HandlerFn func(http.ResponseWriter, *http.Request) error

func (h HandlerFn) Handle(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

type Endpoint struct {
	Handler

	RequestLogger
	Authenticator
	Authorizer
	ErrorPresenter
}

func (fn *Endpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	var ctx context.Context

	// Log all requests, even those that error or panic
	defer func() { fn.LogRequest(req, err) }()

	// Handle panic
	defer fn.stayCalm(w)

	// Authenticate the request
	if ctx, err = fn.Authenticate(req); err != nil {
		fn.PresentError(w, err)
		return
	}
	req = req.WithContext(ctx)

	// Authorize the request
	if err = fn.Authorize(req); err != nil {
		fn.PresentError(w, err)
		return
	}

	// Handle request
	if err = fn.Handle(w, req); err != nil {
		fn.PresentError(w, err)
		return
	}

	// TODO: Figure out per request logging
	// TODO: Implement AARP
}

func (fn *Endpoint) LogRequest(r *http.Request, err error) {
	if fn.RequestLogger == nil {
		return
	}

	fn.RequestLogger.LogRequest(r, err)
}

func (fn *Endpoint) PresentError(w http.ResponseWriter, err error) {
	if fn.ErrorPresenter == nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": %q}`, err.Error()), http.StatusBadRequest)
		return
	}

	fn.ErrorPresenter.PresentError(w, err)
}

func (fn *Endpoint) Authenticate(req *http.Request) (context.Context, error) {
	if fn.Authenticator == nil {
		return req.Context(), nil
	}

	return fn.Authenticator.Authenticate(req)
}

func (fn *Endpoint) Authorize(req *http.Request) error {
	if fn.Authorizer == nil {
		return nil
	}

	return fn.Authorizer.Authorize(req)
}

func (fn *Endpoint) stayCalm(w http.ResponseWriter) {
	if recovered := recover(); recovered != nil {
		switch err := recovered.(type) {
		case error:
			fn.PresentError(w, err)
		default:
			fn.PresentError(w, fmt.Errorf("panic: %v", err))
		}
	}
}
