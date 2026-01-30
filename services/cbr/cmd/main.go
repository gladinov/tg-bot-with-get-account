package main

import (
	"cbr/internal/configs"
	"cbr/internal/handlers"
	"cbr/internal/service"
	"context"
	"log/slog"
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

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("host", conf.CbrHost),
		slog.String("cbr_app_host", conf.Clients.CbrAppApiClient.Host),
		slog.String("cbr_app_port", conf.Clients.CbrAppApiClient.Port))

	logg.Info("initialize service.Transport")
	transport := service.NewTransport(logg, conf.CbrHost)
	logg.Info("initialize service client")
	client := service.NewClient(logg, transport)
	logg.Info("initialize service")
	service := service.NewService(logg, client)
	logg.Info("initialize handlers")
	handlers := handlers.NewHandlers(logg, service)

	logg.Info("initialize router echo")
	router := echo.New()
	router.Use(middleware.CORS())
	router.Use(handlers.ContextHeaderTraceIdMiddleWare)
	router.Use(handlers.LoggerMiddleWare)

	router.POST("/cbr/currencies", handlers.GetAllCurrencies)

	go func() {
		<-ctx.Done()
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	}()
	address := conf.Clients.CbrAppApiClient.GetCbrAppServer()
	logg.Info("run CBR App", slog.String("address", address))
	router.Start(address)
}
