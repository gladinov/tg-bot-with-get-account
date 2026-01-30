package main

import (
	"cbr/internal/clients/cbr"
	"cbr/internal/configs"
	"cbr/internal/handlers"
	"cbr/internal/service"
	"cbr/internal/utils"
	"context"
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

	go func() {
		<-ctx.Done()

		logg.Info("shutting down http server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := router.Shutdown(shutdownCtx); err != nil {
			logg.Error("failed to shutdown server", slog.Any("error", err))
		}
	}()
	address := conf.Clients.CbrAppApiClient.GetCbrAppServer()
	logg.Info("run CBR App", slog.String("address", address))
	if err := router.Start(address); err != nil && err != http.ErrServerClosed {
		logg.Error("server stopped with error", slog.Any("error", err))
	}
}
