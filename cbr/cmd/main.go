package main

import (
	"cbr/internal/configs"
	"cbr/internal/handlers"
	loggerhandler "cbr/internal/handlers/logger"
	"cbr/internal/service"
	"cbr/lib/logger/sl"
	"cbr/pkg/app"
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app.MustInitialize()
	rootPath := app.MustGetRoot()

	cnfgs := configs.MustInitConfig(rootPath)

	logg := sl.NewLogger(cnfgs.Env)

	logg.Info("start app",
		slog.String("env", cnfgs.Env),
		slog.String("host", cnfgs.CbrHost),
		slog.String("server", cnfgs.Http_server))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	router := echo.New()
	transport := service.NewTransport(cnfgs.CbrHost, logg)
	client := service.NewClient(transport, logg)
	srvc := service.NewService(client, logg)
	handlers := handlers.NewHandlers(srvc, logg)

	router.Use(middleware.CORS())
	router.Use(loggerhandler.LoggerMiddleware(logg))

	router.POST("/cbr/currencies", handlers.GetAllCurrencies)

	go func() {
		<-ctx.Done()
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

	}()
	router.Start(cnfgs.Http_server)
}
