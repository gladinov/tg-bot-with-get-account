package main

import (
	"cbr/internal/clients/cbr"
	"cbr/internal/configs"
	"cbr/internal/handlers"
	"cbr/internal/service"
	"cbr/internal/utils"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	sl "github.com/gladinov/mylogger"
	"github.com/gladinov/traceidgenerator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	_ = traceidgenerator.Must()

	conf := configs.MustInitConfig()

	timeLocation := utils.MustGetMoscowLocation()

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("host", conf.CbrHost),
		slog.String("cbr_app_host", conf.Clients.CbrAppApiClient.Host),
		slog.String("cbr_app_port", conf.Clients.CbrAppApiClient.Port))

	logg.Info("initialize cbr.Transport")
	transport := cbr.NewTransport(logg, conf.CbrHost)
	logg.Info("initialize cbr client")
	client := cbr.NewClient(logg, transport)
	logg.Info("initialize service")
	service := service.NewService(logg, client, timeLocation)
	logg.Info("initialize handlers")
	handler := handlers.NewHandlers(logg, service)

	logg.Info("initialize router echo")
	router := echo.New()
	router.Use(middleware.CORS())
	router.Use(handler.ContextHeaderTraceIdMiddleWare)
	router.Use(handler.LoggerMiddleWare)
	router.HTTPErrorHandler = handlers.HTTPErrorHandler(logg)

	router.POST("/cbr/currencies", handler.GetAllCurrencies)
	address := conf.Clients.CbrAppApiClient.GetCbrAppServer()

	httpSrv := &http.Server{
		Addr:         address,
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		logg.Info("run cbr server", slog.String("address", address))
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
