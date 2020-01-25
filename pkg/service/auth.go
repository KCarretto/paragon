package middleware

import (
	"net/http"

	"github.com/kcarretto/paragon/pkg/auth"
)

type Authenticator interface {
	Authenticate(*http.Request) (*auth.Context, error)
}

func WithAuth(svc Authenticator) func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			ctx, err := svc.Authenticate(req)
			if err != nil {
				http.Redirect(w, req, auth.LoginURL, http.StatusTemporaryRedirect)
				return
			}
			req = req.WithContext(ctx)
			h(w, req)
		}
	}
}
