// +build dev

package main

import (
	"net/http"
	"os"
	"strings"
)

func serve(addr string, svc http.Handler) error {
	if disableCORS := os.Getenv("PG_DISABLE_CORS"); disableCORS == "" {
		return http.ListenAndServe(addr, svc)
	}

	withCORSDisabled := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.RequestURI == "/" || r.RequestURI == "/login" || strings.HasPrefix(r.RequestURI, "/app") {
			http.Redirect(w, r, "http://127.0.0.1:8080"+r.URL.Path, http.StatusTemporaryRedirect)
		}
		svc.ServeHTTP(w, r)

	}

	return http.ListenAndServe(addr, http.HandlerFunc(withCORSDisabled))
}
