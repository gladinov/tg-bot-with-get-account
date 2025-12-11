package handlers

import (
	"cbr/internal/service"
	"cbr/lib/models"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	service service.CurrencyService
	logger  *slog.Logger
}

func NewHandlers(srvc service.CurrencyService, logger *slog.Logger) *Handlers {
	return &Handlers{service: srvc,
		logger: logger}
}

func (h *Handlers) GetAllCurrencies(c echo.Context) error {
	const op = "handlers.GetAllCurrencies"
	logg := h.logger.With(
		slog.String("function", op),
	)
	logg.Info("start " + op)

	var currencyRequest CurrencyRequest

	logg.Debug("bind input requst")

	err := c.Bind(&currencyRequest)
	if err != nil {
		logg.Error("echo.Bind error", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	logg.Debug("service.GetAllCurrencies")
	currencies, err := h.service.GetAllCurrencies(currencyRequest.Date)
	if err != nil {
		logg.Error("service.GetAllCurrencies error", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get currencies"})
	}
	logg.Info("end " + op)
	requestID := c.Request().Header.Get(models.RequestIDHeader)
	logg.Debug(requestID)
	c.Response().Header().Set(models.RequestIDHeader, requestID)
	return c.JSON(http.StatusOK, currencies)
}
