package handlers

import (
	"bytes"
	"cbr/internal/service"
	"cbr/internal/service/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAllCurrencies(t *testing.T) {

	cases := []struct {
		name           string
		request        any
		setupMock      func(*mocks.CurrencyService)
		expectedStatus int
		expectedBody   string
		assertNoCalls  bool
	}{
		{
			name:    "Succses",
			request: CurrencyRequest{Date: time.Date(2025, time.November, 10, 0, 0, 0, 0, time.UTC)},
			setupMock: func(mockCurrencyService *mocks.CurrencyService) {
				mockCurrencyService.On("GetAllCurrencies", mock.AnythingOfType("time.Time")).Once().
					Return(service.HappyPathCurrencies, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   service.HappyPathCurrenciesInBytes,
			assertNoCalls:  false,
		},
		{
			name:           "Err: Invalid Request",
			request:        `{"invalid json"}`,
			setupMock:      func(mockCurrencyService *mocks.CurrencyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid request"}`,
			assertNoCalls:  true,
		},
		{
			name:    "Err: GetAllCurencies err",
			request: CurrencyRequest{Date: time.Date(2025, time.November, 10, 0, 0, 0, 0, time.UTC)},
			setupMock: func(mockCurrencyService *mocks.CurrencyService) {
				mockCurrencyService.On("GetAllCurrencies", mock.AnythingOfType("time.Time")).Once().
					Return(service.CurrenciesResponce{}, errors.New("op: service.GetAllCurrencies, error: failed to load Moscow location"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "could not get currencies"}`,
			assertNoCalls:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			var jsonData []byte
			switch body := tc.request.(type) {
			case []byte:
				jsonData = body
			case CurrencyRequest:
				jsonData, _ = json.Marshal(body)

			case string:
				jsonData = []byte(body)
			}
			mockService := mocks.NewCurrencyService(t)
			tc.setupMock(mockService)
			ctx, rec := createTestContext(http.MethodPost, "/cbr/currencies", jsonData)

			hndlrs := NewHandlers(mockService)
			err := hndlrs.GetAllCurrencies(ctx)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				mockService.AssertNotCalled(t, "GetAllCurrencies")

			} else {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func createTestContext(method, path string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))

	e := echo.New()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
