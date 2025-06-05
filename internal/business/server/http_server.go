package server

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/samber/oops"

	slogctx "github.com/veqryn/slog-context"

	"github.com/openkcm/dummy-service/internal/config"
)

// registerHandlers registers the default http handlers for the status server
func registerHandlers(mux *http.ServeMux, cfg *config.Config) {
	mux.HandleFunc("/ping", pingHandlerFunc(cfg))
}

// createStatusServer creates a status http server using the given config
func createHTTPServer(ctx context.Context, cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	registerHandlers(mux, cfg)

	slogctx.Info(ctx, "Creating HTTP server", "address", cfg.HTTP.Address)

	return &http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: mux,
	}
}

// StartHTTPServer starts the gRPC server using the given config.
func StartHTTPServer(ctx context.Context, cfg *config.Config) error {
	if err := initMeters(ctx, cfg); err != nil {
		return err
	}

	server := createHTTPServer(ctx, cfg)

	slogctx.Info(ctx, "Starting HTTP listener", "address", server.Addr)

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return oops.In("HTTP Server").
			WithContext(ctx).
			Wrapf(err, "Failed creating HTTP listener")
	}

	go func() {
		slogctx.Info(ctx, "Starting HTTP server", "address", server.Addr)

		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogctx.Error(ctx, "ErrorField serving HTTP endpoint", "error", err)
		}

		slogctx.Info(ctx, "Stopped HTTP server")
	}()

	<-ctx.Done()

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, cfg.HTTP.ShutdownTimeout)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return oops.In("HTTP Server").
			WithContext(ctx).
			Wrapf(err, "Failed shutting down HTTP server")
	}

	slogctx.Info(ctx, "Completed graceful shutdown of HTTP server")

	return nil
}
