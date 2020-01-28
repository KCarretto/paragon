// +build !dev

package main

import (
	"net/http"
)

func serve(addr string, svc http.Handler) error {
	return http.ListenAndServe(addr, svc)
}
