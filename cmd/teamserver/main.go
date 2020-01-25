package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kcarretto/paragon/pkg/middleware"
	"github.com/kcarretto/paragon/pkg/teamserver"
	"go.uber.org/zap"
)

// HTTPAddr for the teamserver to listen on. Configure with HTTP_ADDR env variable.
var (
	HTTPAddr = "127.0.0.1:8080"
	start    time.Time
)

// DefaultTopic is the default pubsub topic where events will be published. It can be configured
// by setting the PUB_TOPIC environment variable.
const DefaultTopic = "TOPIC"

func init() {
	start = time.Now()
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		HTTPAddr = addr
	}
}

// ServeHTTP handles http requests for the Teamserver.
func ServeHTTP(logger *zap.Logger, addr string, svc *teamserver.Service) error {
	router := http.NewServeMux()
	svc.HTTP(router)

	defer func() {
		logger.Info("Teamserver Stopped",
			zap.Duration("uptime", time.Since(start)),
		)
	}()

	logger.Info("Teamserver Started",
		zap.String("listen_on", addr),
		zap.Duration("start_latency", time.Since(start)),
	)

	return http.ListenAndServe(
		addr,
		middleware.Chain(
			router,
			middleware.WithPanicHandling,
			middleware.WithLogging(logger),
		))
}

func main() {
	ctx := context.Background()
	logger := newLogger().Named("svc.teamserver")
	logger.Debug("Initializing teamserver", zap.Time("start", start))

	topic := os.Getenv("PUB_TOPIC")
	if topic == "" {
		log.Println("[WARN] No PUB_TOPIC environment variable set to publish events, is this a mistake?")
	}
	publisher := newPublisher(ctx, topic)

	graph := newGraph(ctx)
	defer graph.Close()

	svc := &teamserver.Service{
		Log:    logger,
		Graph:  graph,
		Events: publisher,
	}

	ServeHTTP(logger.Named("http"), HTTPAddr, svc)
}
