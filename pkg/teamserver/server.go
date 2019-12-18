package teamserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kcarretto/paragon/api/ent/credentials"
	"github.com/kcarretto/paragon/api/ent/job_templates"
	"github.com/kcarretto/paragon/api/ent/jobs"
	"github.com/kcarretto/paragon/api/ent/tags"
	"github.com/kcarretto/paragon/api/ent/targets"
	"github.com/kcarretto/paragon/api/ent/tasks"
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/pkg/middleware"
)

// Server handles c2 messages and replies with new tasks for the c2 to send out.
type Server struct {
	Log       *zap.Logger
	EntClient *ent.Client
}

func (srv *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status": "OK",
	}
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to marshal the json for the status", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response data: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

// Run begins the handlers for processing the subscriptions to the `tasks.claimed` and `tasks.executed` topics
func (srv *Server) Run() {
	ctx := context.Background()
	svcRouter := runtime.NewServeMux()

	targetSVC := &targets.Service{
		EntClient: srv.EntClient,
	}
	if err := targets.RegisterTargetsHandlerServer(ctx, svcRouter, targetSVC); err != nil {
		panic(err)
	}

	taskSVC := &tasks.Service{
		EntClient: srv.EntClient,
	}
	if err := tasks.RegisterTasksHandlerServer(ctx, svcRouter, taskSVC); err != nil {
		panic(err)
	}

	tagSVC := &tags.Service{
		EntClient: srv.EntClient,
	}
	if err := tags.RegisterTagsHandlerServer(ctx, svcRouter, tagSVC); err != nil {
		panic(err)
	}

	jobSVC := &jobs.Service{
		EntClient: srv.EntClient,
	}
	if err := jobs.RegisterJobsHandlerServer(ctx, svcRouter, jobSVC); err != nil {
		panic(err)
	}

	jobTemplateSVC := &job_templates.Service{
		EntClient: srv.EntClient,
	}
	if err := job_templates.RegisterJobTemplatesHandlerServer(ctx, svcRouter, jobTemplateSVC); err != nil {
		panic(err)
	}

	credentialSVC := &credentials.Service{
		EntClient: srv.EntClient,
	}
	if err := credentials.RegisterCredentialsHandlerServer(ctx, svcRouter, credentialSVC); err != nil {
		panic(err)
	}

	router := http.NewServeMux()

	router.Handle("/api/v1/", svcRouter)
	router.HandleFunc("/status", srv.handleStatus)

	if err := http.ListenAndServe("0.0.0.0:80", middleware.Chain(router, middleware.WithPanicHandling)); err != nil {
		panic(err)
	}
}
