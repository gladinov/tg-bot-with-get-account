package main

import (
	"cbr/internal/configs"
	"cbr/internal/handlers"
	"cbr/internal/service"
	"cbr/pkg/app"
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const ()

func main() {
	app.MustInitialize()
	rootPath := app.MustGetRoot()

	cnfgs := configs.MustInitConfig(rootPath)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	router := echo.New()
	transport := service.NewTransport(cnfgs.CbrHost)
	client := service.NewClient(transport)
	srvc := service.NewService(client)
	handlers := handlers.NewHandlers(srvc)

	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	router.POST("/cbr/currencies", handlers.GetAllCurrencies)

	go func() {
		<-ctx.Done()
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

	}()
	router.Start(cnfgs.Http_server)
}
