package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type Middleware struct {
	Middleware []func(http.Handler) http.HandlerFunc
}


func Wrap(h http.Handler) *Endpoint {
	return &Endpoint{
		Handler: h,
	}
}

func (e *Endpoint) WithMiddleware(middleware ...func(http.Handler) http.HandlerFunc) *Endpoint {
	e.Middleware = append(e.Middleware, middleware...)
	return e
}

func (e *Endpoint) Auth(svc Authenticator) *Endpoint {
	return e.WithMiddleware(WithAuth(svc))
}

func (e *Endpoint) Logging(logger *zap.Logger) *Endpoint {
	return e.WithMiddleware(WithLogging(logger))
}

func (e *Endpoint) Wrap(h http.Handler) http.HandlerFunc){

}
