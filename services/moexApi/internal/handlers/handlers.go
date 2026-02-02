package handlers

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/models"
	"main/internal/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	errGetData            error = errors.New("could not get data")
	errInvalidRequestBody error = errors.New("invalid request body")
)

type Handlers struct {
	logger  *slog.Logger
	service service.ServiceClient
}

func NewHandlers(logger *slog.Logger, service service.ServiceClient) *Handlers {
	return &Handlers{
		logger:  logger,
		service: service,
	}
}

func (h *Handlers) GetSpecifications(c echo.Context) error {
	const op = "handlers.GetSpecifications"
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logg := h.logger.With(slog.String("op", op))
	logg.DebugContext(ctx, "start")

	var req models.SpecificationsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidRequestBody)
	}
	resp, err := h.service.GetSpecifications(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetData)
	}
	return c.JSON(http.StatusOK, resp)
}
