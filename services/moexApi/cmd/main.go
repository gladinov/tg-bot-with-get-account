package main

import (
	"context"
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

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := router.Shutdown(shutdownCtx); err != nil {
			logg.Error("Failed to shutdown server:", slog.String("error", err.Error()))
		}
	}()
	address := conf.Clients.MoexApiAppClient.GetMoexApiAppClientAddress()
	logg.Info("run MOEX API App", slog.String("address", address))

	if err := router.Start(address); err != nil && err != http.ErrServerClosed {
		logg.Error("server start failed", slog.Any("error", err))
	}
}
