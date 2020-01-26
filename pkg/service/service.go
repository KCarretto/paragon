package service

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

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

	Log *zap.Logger
	Authenticator
	Authorizer
	ErrorPresenter
}

func (fn *Endpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	var ctx context.Context

	logger := fn.logger(req)

	// Log all requests, even those that error or panic
	defer func() {
		logger.Info("Request Handled", zap.Error(err))
	}()

	// Handle panic
	defer fn.stayCalm(w, logger)

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

func (fn *Endpoint) logger(req *http.Request) *zap.Logger {
	if fn.Log == nil {
		fn.Log = zap.NewNop()
	}

	// TODO: Add request tracing info
	return fn.Log.With(
		zap.String("method", req.Method),
		zap.String("uri", req.RequestURI),
		zap.Int64("length", req.ContentLength),
		zap.String("remote_addr", req.RemoteAddr),
		zap.String("user_agent", req.UserAgent()),
	)
}

func (fn *Endpoint) stayCalm(w http.ResponseWriter, logger *zap.Logger) {
	if recovered := recover(); recovered != nil {
		logger.Error("Request resulted in panic", zap.Any("error", recovered))

		switch err := recovered.(type) {
		case error:
			fn.PresentError(w, err)
		default:
			fn.PresentError(w, fmt.Errorf("panic: %v", err))
		}
	}
}
