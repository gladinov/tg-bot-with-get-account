package main

import (
	"context"
	"log/slog"
	"main/internal/configs"
	"main/internal/handlers"
	"main/internal/service"
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

	conf := configs.MustLoad()

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("host", conf.MoexHost),
		slog.String("cbr_app_host", conf.Clients.MoexApiAppClient.Host),
		slog.String("cbr_app_port", conf.Clients.MoexApiAppClient.Port))

	logg.Info("initialize SpecificationService")
	service := service.NewSpecificationService(logg, conf.MoexHost)

	logg.Info("initialize handlers")
	handlers := handlers.NewHandlers(logg, service)

	logg.Info("initialize router echo")
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(handlers.ContextHeaderTraceIdMiddleWare)
	e.Use(handlers.LoggerMiddleWare)

	e.POST("/moex/specifications", handlers.GetSpecifications)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(shutdownCtx); err != nil {
			logg.Error("Failed to shutdown server:", slog.String("error", err.Error()))
		}
	}()
	address := conf.Clients.MoexApiAppClient.GetMoexApiAppClientAddress()
	logg.Info("run MOEX API App", slog.String("address", address))
	e.Start(address)
}
