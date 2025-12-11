package main

import (
	"main/internal/configs"
	"main/internal/handlers"
	"main/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//app.MustInitialize()
	//rootPath := app.MustGetRoot()

	config := configs.MustLoad()

	service := service.NewSpecificationService(config.MoexHost)
	handlers := handlers.NewHandlers(service)
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/moex/specifications", handlers.GetSpecifications)

	e.Start(config.Clients.MoexApiAppClient.GetMoexApiAppClientAddress())
}
