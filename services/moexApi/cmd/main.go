package main

import (
	"context"
	"errors"
	"log/slog"
	"moex/internal/clients/moex"
	"moex/internal/configs"
	"moex/internal/handlers"
	"moex/internal/service"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	sl "github.com/gladinov/mylogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	conf := configs.MustLoad()

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("host", conf.MoexHost),
		slog.String("cbr_app_host", conf.Clients.MoexApiAppClient.Host),
		slog.String("cbr_app_port", conf.Clients.MoexApiAppClient.Port))

	logg.Info("initialize Transport")
	transport := moex.NewTransport(logg, conf.MoexHost)
	logg.Info("initialize client")
	moexClient := moex.NewMoexClient(logg, transport)
	logg.Info("initialize service")
	service := service.NewServiceClient(logg, moexClient)

	logg.Info("initialize handlers")
	handler := handlers.NewHandlers(logg, service)

	logg.Info("initialize router echo")
	router := echo.New()

	router.Use(middleware.CORS())
	router.Use(handler.ContextHeaderTraceIdMiddleWare)
	router.Use(handler.LoggerMiddleWare)
	router.HTTPErrorHandler = handlers.HTTPErrorHandler(logg)

	router.POST("/moex/specifications", handler.GetSpecifications)

	address := conf.Clients.MoexApiAppClient.GetMoexApiAppClientAddress()

	httpSrv := &http.Server{
		Addr:         address,
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		logg.Info("run MOEX API App", slog.String("address", address))
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		logg.InfoContext(ctx, "Shutdown signal received")
	case err := <-errCh:
		logg.ErrorContext(ctx, "server stopped with error", slog.Any("error", err))
	}
	gracefulShutdown(ctx, logg, httpSrv)
}

func gracefulShutdown(ctx context.Context, logg *slog.Logger, httpSrv *http.Server) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logg.ErrorContext(ctx, "Forced shutdown", slog.Any("err", err))
	}
	logg.InfoContext(ctx, "Server exited gracefully")
}
