package www

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Service provides HTTP handlers for serving the Teamserver UI.
type Service struct {
	Log *zap.Logger
}

func (svc Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	f, err := App.Open("index.html")
	if err != nil {
		http.Error(w, "Failed to load index.html", http.StatusNotFound)
		return
	}

	var modtime time.Time
	if info, err := f.Stat(); err == nil && info != nil {
		modtime = info.ModTime()
	}

	http.ServeContent(w, r, "index.html", modtime, f)
}

// HTTP registers http handlers for the Teamserver UI.
func (svc *Service) HTTP(router *http.ServeMux) {
	router.Handle("/app/", http.StripPrefix("/app", http.FileServer(App)))
	router.HandleFunc("/", svc.handleIndex)
}
