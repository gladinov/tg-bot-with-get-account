//go:build unit

package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"moex/internal/models"
	"moex/internal/service/mocks"
	"moex/internal/testdata/factories"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSpecifications(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockService := mocks.NewServiceClient(t)
	h := NewHandlers(logger, mockService)
	e := echo.New()
	e.HTTPErrorHandler = HTTPErrorHandler(logger)

	t.Run("Success: returns expected specification", func(t *testing.T) {
		expectedStatus := http.StatusOK

		want := factories.NewValues()
		inputBody := map[string]any{
			"ticker": "RU000A10B9Q9",
			"date":   "2025-11-16T15:12:46.3365285+03:00",
		}

		mockService.On(
			"GetSpecifications",
			mock.Anything,
			mock.MatchedBy(func(req models.SpecificationsRequest) bool {
				return req.Ticker == "RU000A10B9Q9" &&
					!req.Date.IsZero()
			}),
		).Return(want, nil).Once()

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/moex/specifications", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/moex/specifications", h.GetSpecifications)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		responseBody := rec.Body.String()
		assert.Contains(t, responseBody, want.ShortName.Value)
		assert.Contains(t, responseBody, want.TradeDate.Value)
		assert.Contains(t, responseBody, want.MaturityDate.Value)
		assert.Contains(t, responseBody, want.OfferDate.Value)
		assert.Contains(t, responseBody, want.BuybackDate.Value)

		mockService.AssertExpectations(t)
	})

	t.Run("Err: Bind", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest

		wantMessage := errInvalidRequestBody.Error()
		inputBody := map[string]any{
			"ticker": 12,
			"date":   "2025-11-16T15:12:46.3365285+03:00",
		}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/moex/specifications", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/moex/specifications", h.GetSpecifications)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		responseBody := rec.Body.String()

		assert.Contains(t, responseBody, wantMessage)
	})

	t.Run("Err: validateRequest ticker err", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest

		wantMessage := errInvalidRequestBody.Error()
		inputBody := map[string]any{
			"ticker": "",
			"date":   "2025-11-16T15:12:46.3365285+03:00",
		}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/moex/specifications", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/moex/specifications", h.GetSpecifications)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		responseBody := rec.Body.String()

		assert.Contains(t, responseBody, wantMessage)
	})

	t.Run("Err: validateRequest: empty ticker err", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest

		wantMessage := errInvalidRequestBody.Error()
		inputBody := map[string]any{
			"date": "2025-11-16T15:12:46.3365285+03:00",
		}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/moex/specifications", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/moex/specifications", h.GetSpecifications)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		responseBody := rec.Body.String()

		assert.Contains(t, responseBody, wantMessage)
	})

	t.Run("Err: validateRequest zero date err", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest

		wantMessage := errInvalidRequestBody.Error()
		inputBody := map[string]any{
			"ticker": "RU000A10B9Q9",
			"date":   "0001-01-01T00:00:00Z",
		}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/moex/specifications", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/moex/specifications", h.GetSpecifications)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		responseBody := rec.Body.String()

		assert.Contains(t, responseBody, wantMessage)
	})

	t.Run("Err: validateRequest empty date err", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest

		wantMessage := errInvalidRequestBody.Error()
		inputBody := map[string]any{
			"ticker": "RU000A10B9Q9",
		}

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/moex/specifications", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/moex/specifications", h.GetSpecifications)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)

		responseBody := rec.Body.String()

		assert.Contains(t, responseBody, wantMessage)
	})

	t.Run("Err: GetSpecifications err ", func(t *testing.T) {
		expectedStatus := http.StatusInternalServerError

		inputBody := map[string]any{
			"ticker": "RU000A10B9Q9",
			"date":   "2025-11-16T15:12:46.3365285+03:00",
		}

		mockService.On(
			"GetSpecifications",
			mock.Anything,
			mock.MatchedBy(func(req models.SpecificationsRequest) bool {
				return req.Ticker == "RU000A10B9Q9" &&
					!req.Date.IsZero()
			}),
		).Return(models.Values{}, errors.New("could not get specification from moexClient")).Once()

		bodyBytes, _ := json.Marshal(inputBody)
		req := httptest.NewRequest(http.MethodPost, "/moex/specifications", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.POST("/moex/specifications", h.GetSpecifications)
		e.ServeHTTP(rec, req)

		assert.Equal(t, expectedStatus, rec.Code)
		assert.Contains(t, rec.Body.String(), errGetData.Error())

		mockService.AssertExpectations(t)
	})
}
