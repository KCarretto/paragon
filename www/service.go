package www

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kcarretto/paragon/pkg/service"
	"go.uber.org/zap"
)

// Service provides HTTP handlers for serving the Teamserver UI.
type Service struct {
	Log *zap.Logger
}

func (svc Service) HandleIndex(w http.ResponseWriter, r *http.Request) error {
	f, err := App.Open("index.html")
	if err != nil {
		return fmt.Errorf("Failed to load index.html")
	}

	var modtime time.Time
	if info, err := f.Stat(); err == nil && info != nil {
		modtime = info.ModTime()
	}

	http.ServeContent(w, r, "index.html", modtime, f)
	return nil
}

// HTTP registers http handlers for the Teamserver UI.
func (svc *Service) HTTP(router *http.ServeMux) {
	app := &service.Endpoint{
		Handler: service.HTTPHandler(http.FileServer(App)),
	}
	index := &service.Endpoint{
		Handler: service.HandlerFn(svc.HandleIndex),
	}
	router.Handle("/app/", http.StripPrefix("/app", app))
	router.Handle("/", index)
}
