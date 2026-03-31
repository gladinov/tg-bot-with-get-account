//go:build unit

package handlers

import (
	"bytes"
	"cbr/internal/models"
	"cbr/internal/service/mocks"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllCurrencies(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockService := mocks.NewCurrencyService(t)
	h := NewHandlers(logger, mockService)
	e := echo.New()
	e.HTTPErrorHandler = HTTPErrorHandler(logger)

	t.Run("sucsess", func(t *testing.T) {
		expectedStatus := http.StatusOK
		expectedBody := models.CurrenciesResponce{}
		inputBody := map[string]any{"date": "2025-11-16T15:12:46.3365285+03:00"}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/cbr/currencies", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		mockService.On("GetAllCurrencies", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(expectedBody, nil).Once()

		e.POST("/cbr/currencies", h.GetAllCurrencies)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		var respBody models.CurrenciesResponce
		err := json.Unmarshal(rec.Body.Bytes(), &respBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, respBody)

		mockService.AssertExpectations(t)
	})
	t.Run("invalid request", func(t *testing.T) {
		expectedBody := errInvalidRequestBody.Error()
		expectedStatus := http.StatusBadRequest
		inputBody := map[string]any{"date": 12}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/cbr/currencies", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/cbr/currencies", h.GetAllCurrencies)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		var respBody map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &respBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, respBody["error"])
	})
	t.Run("GetAllCurrencies err", func(t *testing.T) {
		expectedBody := errGetData.Error()
		expectedStatus := http.StatusInternalServerError
		inputBody := map[string]any{"date": "2025-11-16T15:12:46.3365285+03:00"}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/cbr/currencies", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		mockService.On("GetAllCurrencies", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(models.CurrenciesResponce{}, errors.New("failed to get all currencies from client")).Once()

		e.POST("/cbr/currencies", h.GetAllCurrencies)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		var respBody map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &respBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, respBody["error"])
	})
}
