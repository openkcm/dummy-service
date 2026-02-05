package server

import (
	"context"
	"net"

	"github.com/openkcm/common-sdk/pkg/health"
	"github.com/openkcm/common-sdk/pkg/otlp"
	"github.com/samber/oops"
	"google.golang.org/grpc"

	slogctx "github.com/veqryn/slog-context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/openkcm/dummy-service/internal/config"
)

// createGRPCServer creates the gRPC server using the given config
func createGRPCServer(_ context.Context, _ *config.Config) (*grpc.Server, error) {
	// Create the gRPC server options
	opts := make([]grpc.ServerOption, 0, 1)
	opts = append(opts, grpc.StatsHandler(otlp.NewServerHandler()))

	grpcServer := grpc.NewServer(opts...)

	// Create the gRPC server
	return grpcServer, nil
}

// StartGRPCServer starts the gRPC server using the given config.
func StartGRPCServer(ctx context.Context, cfg *config.Config) error {
	// Create the gRPC server
	grpcServer, err := createGRPCServer(ctx, cfg)
	if err != nil {
		return oops.In("GRPC Server").
			WithContext(ctx).
			Wrapf(err, "Failed creation of gRPC server")
	}

	// Register the servers with the gRPC server
	//envoy_auth.RegisterAuthorizationServer(grpcServer, grpcServer)
	healthpb.RegisterHealthServer(grpcServer, &health.GRPCServer{})

	// Start the listener
	listener, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return oops.In("GRPC Server").
			WithContext(ctx).
			Wrapf(err, "Failed to create the listener")
	}

	// Serve the gRPC server
	go func() {
		slogctx.Info(ctx, "Starting gRPC Server", "address", cfg.GRPC.Address)

		err = grpcServer.Serve(listener)

		if err != nil {
			slogctx.Error(ctx, "ErrorField serving gRPC endpoint", "error", err)
		}

		slogctx.Info(ctx, "Stopped gRPC server")
	}()

	<-ctx.Done()

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, cfg.GRPC.ShutdownTimeout)
	defer shutdownRelease()

	grpcServer.GracefulStop()
	slogctx.Info(shutdownCtx, "Completed graceful shutdown of gRPC server")

	return nil
}
