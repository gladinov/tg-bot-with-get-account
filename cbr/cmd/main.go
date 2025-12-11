package main

import (
	"cbr/internal/configs"
	"cbr/internal/handlers"
	loggerhandler "cbr/internal/handlers/logger"
	"cbr/internal/service"
	"cbr/lib/logger/sl"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//for local
	//app.MustInitialize()
	//rootPath := app.MustGetRoot()
	rootPath := os.Getenv("ROOT_PATH")
	if rootPath == "" {
		panic("ROOT_PATH environment variable is required")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH environment variable is required")
	}

	cnfgs := configs.MustInitConfig(rootPath, configPath)

	logg := sl.NewLogger(cnfgs.Env)

	logg.Info("start app",
		slog.String("env", cnfgs.Env),
		slog.String("host", cnfgs.CbrHost),
		slog.String("cbr_app_host", cnfgs.CbrAppHost),
		slog.String("cbr_app_port", strconv.Itoa(cnfgs.CbrAppPort)))

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
	router.Start(cnfgs.GetCbrAppServer())
}
