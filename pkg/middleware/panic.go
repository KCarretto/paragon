package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WithPanicHandling is used to handle all panics on the handlers and pass back an error
func WithPanicHandling(f http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if fatalErr := recover(); fatalErr != nil {
				err := fmt.Errorf("Uh oh: %v", fatalErr)
				ret := map[string]interface{}{
					"error":   err.Error(),
					"code":    2,
					"message": err.Error(),
				}
				jsonBytes, _ := json.Marshal(ret)
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write(jsonBytes)
			}
		}()
		f.ServeHTTP(rw, r)
	}
}
