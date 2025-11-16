package handlers

import (
	"cbr/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	service service.CurrencyService
}

func NewHandlers(srvc service.CurrencyService) *Handlers {
	return &Handlers{service: srvc}
}

func (h *Handlers) GetAllCurrencies(c echo.Context) error {
	var currencyRequest CurrencyRequest
	err := c.Bind(&currencyRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	currencies, err := h.service.GetAllCurrencies(currencyRequest.Date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get currencies"})
	}

	return c.JSON(http.StatusOK, currencies)
}
