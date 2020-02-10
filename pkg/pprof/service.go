package pprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/kcarretto/paragon/pkg/service"
	"go.uber.org/zap"
)

// Service provides HTTP handlers that provide profiling information about the service.
type Service struct {
	Log   *zap.Logger
	Auth  service.Authenticator
	Authz service.Authorizer
}

// HTTP registers http handlers that serve profiling information.
func (svc *Service) HTTP(router *http.ServeMux) {
	index := &service.Endpoint{
		Log:           svc.Log.Named("index"),
		Authenticator: svc.Auth,
		Authorizer:    svc.Authz,
		Handler:       service.HTTPHandler(http.HandlerFunc(pprof.Index)),
	}
	cmdline := &service.Endpoint{
		Log:           svc.Log.Named("cmdline"),
		Authenticator: svc.Auth,
		Authorizer:    svc.Authz,
		Handler:       service.HTTPHandler(http.HandlerFunc(pprof.Cmdline)),
	}
	profile := &service.Endpoint{
		Log:           svc.Log.Named("profile"),
		Authenticator: svc.Auth,
		Authorizer:    svc.Authz,
		Handler:       service.HTTPHandler(http.HandlerFunc(pprof.Profile)),
	}
	symbol := &service.Endpoint{
		Log:           svc.Log.Named("symbol"),
		Authenticator: svc.Auth,
		Authorizer:    svc.Authz,
		Handler:       service.HTTPHandler(http.HandlerFunc(pprof.Symbol)),
	}
	trace := &service.Endpoint{
		Log:           svc.Log.Named("trace"),
		Authenticator: svc.Auth,
		Authorizer:    svc.Authz,
		Handler:       service.HTTPHandler(http.HandlerFunc(pprof.Trace)),
	}

	router.Handle("/debug/pprof/", index)
	router.Handle("/debug/pprof/cmdline", cmdline)
	router.Handle("/debug/pprof/profile", profile)
	router.Handle("/debug/pprof/symbol", symbol)
	router.Handle("/debug/pprof/trace", trace)
}
