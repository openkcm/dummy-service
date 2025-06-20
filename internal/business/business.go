package business

import (
	"context"
	"errors"
	"log/slog"

	slogctx "github.com/veqryn/slog-context"

	"github.com/openkcm/dummy-service/internal/business/server"
	"github.com/openkcm/dummy-service/internal/config"
)

func Main(ctx context.Context, cfg *config.Config) error {
	// Application Business Logic
	//------------------------------------------------------------------
	// Example of Slog using Context
	// TODO: To be deleted
	// Add attributes directly to the logger in the context:
	ctx = slogctx.With(ctx, "rootKey", "rootValue")

	// With and wrapper methods have the same args signature as slog methods,
	// and can take a mix of slog.Attr and key-value pairs.
	ctx = slogctx.With(slogctx.WithGroup(ctx, "someGroup"),
		slog.String("email", "n.nicora@sap.com"),
		slog.String("jwt-token", "eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTc0MDkwOTM5NiwiaWF0IjoxNzQwOTA5Mzk2fQ.cxK9EREN6pviQ4E99US60QSQRlX27NgxrUEJLteFa1o"),
		slog.Bool("someBool", true),
	)

	slogctx.Error(ctx, "main message",
		slogctx.Err(errors.New("huraaaaaa")),
		slog.String("mainKey", "mainValue"))
	// Example of Slog using Context
	//------------------------------------------------------------------

	//return server.StartGRPCServer(ctx, cfg)
	return server.StartHTTPServer(ctx, cfg)
}
