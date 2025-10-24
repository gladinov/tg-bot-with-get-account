package handlers

import (
	"main/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	service service.Service
}

func NewHandlers(service service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) GetSpecifications(c echo.Context) error {
	var req service.SpecificationsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	resp, err := h.service.GetSpecifications(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get specifications"})
	}
	return c.JSON(http.StatusOK, resp)
}
