package debug

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func prepareBasicAuth(username, password string) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}

			b, err := base64.StdEncoding.DecodeString(s[1])
			if err != nil {
				http.Error(w, "Malformatted token", http.StatusBadRequest)
				return
			}

			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}

			if pair[0] != username || pair[1] != password {
				http.Error(w, "Not authorized", 401)
				return
			}

			h.ServeHTTP(w, r)
		}
	}
}
