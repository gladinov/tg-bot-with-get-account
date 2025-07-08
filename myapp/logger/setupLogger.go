package loggAdapter

import (
	"log/slog"
	"os"
)

func SetupLogger() *Adapter {

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	slogLogger := New(log)

	return slogLogger
}
