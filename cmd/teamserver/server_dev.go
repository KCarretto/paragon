// +build dev

package main

import (
	"fmt"
	"net/http"
)

func serve(addr string, svc http.Handler) error {
	withCORSDisabled := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		fmt.Printf("REQUEST URI: %v\n")
		fmt.Printf("URL PATH: %v\n", r.URL.Path)

		if r.RequestURI == "/" || r.RequestURI == "/login" {
			http.Redirect(w, r, "http://127.0.0.1:8080"+r.URL.Path, http.StatusTemporaryRedirect)
		}
		svc.ServeHTTP(w, r)

	}

	return http.ListenAndServe(addr, http.HandlerFunc(withCORSDisabled))
}
