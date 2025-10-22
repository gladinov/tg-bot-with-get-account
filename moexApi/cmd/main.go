package main

import (
	"main/internal/handlers"
	"main/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	moexHost = "iss.moex.com"
)

func main() {
	service := service.NewSpecificationService(moexHost)
	handlers := handlers.NewHandlers(service)
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/specifications", handlers.GetSpecifications)

	e.Start("localhost:8081")
}
