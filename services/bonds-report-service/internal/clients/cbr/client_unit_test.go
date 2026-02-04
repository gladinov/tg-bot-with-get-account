//go:build unit

package cbr

import (
	"bonds-report-service/internal/clients/cbr/mocks"
	factories "bonds-report-service/internal/clients/cbr/testdata"
	"bonds-report-service/internal/models"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAllCurrencies(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	ctx := context.Background()
	date := time.Now()

	t.Run("Success", func(t *testing.T) {
		curr := factories.NewCurrency()
		currResp := factories.NewCurrenciesResponce("", curr)
		wantMock := factories.NewHTTPResponse(200, currResp)
		transportMock := mocks.NewTransportClient(t)
		transportMock.On("DoRequest",
			ctx,
			mock.MatchedBy(func(s string) bool { return s != "" }),
			mock.Anything,
			mock.Anything).
			Return(wantMock, nil)
		CbrClient := NewCbrClient(logger, transportMock)
		got, err := CbrClient.GetAllCurrencies(ctx, date)
		require.NoError(t, err)
		require.Len(t, got.Currencies, 1)
		require.Equal(t, currResp.Currencies[0].CharCode, got.Currencies[0].CharCode)
	})

	t.Run("Transport error", func(t *testing.T) {
		transportMock := &mocks.TransportClient{}
		transportMock.On(
			"DoRequest",
			ctx,
			path.Join("cbr", "currencies"),
			url.Values{},
			mock.AnythingOfType("*bytes.Buffer"),
		).Return(nil, errors.New("network failure")).Once()

		client := NewCbrClient(logger, transportMock)
		_, err := client.GetAllCurrencies(ctx, date)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "network failure")

		transportMock.AssertExpectations(t)
	})

	t.Run("HTTP 400 Bad Request", func(t *testing.T) {
		transportMock := &mocks.TransportClient{}
		transportMock.On(
			"DoRequest",
			ctx,
			path.Join("cbr", "currencies"),
			url.Values{},
			mock.AnythingOfType("*bytes.Buffer"),
		).Return(models.NewHTTPResponse(http.StatusBadRequest, []byte(`{}`)), nil).Once()

		client := NewCbrClient(logger, transportMock)
		_, err := client.GetAllCurrencies(ctx, date)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bad request")

		transportMock.AssertExpectations(t)
	})

	t.Run("HTTP 500 Internal Server Error", func(t *testing.T) {
		transportMock := &mocks.TransportClient{}
		transportMock.On(
			"DoRequest",
			ctx,
			path.Join("cbr", "currencies"),
			url.Values{},
			mock.AnythingOfType("*bytes.Buffer"),
		).Return(models.NewHTTPResponse(http.StatusInternalServerError, []byte(`{}`)), nil).Once()

		client := NewCbrClient(logger, transportMock)
		_, err := client.GetAllCurrencies(ctx, date)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "internal server error")

		transportMock.AssertExpectations(t)
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		transportMock := &mocks.TransportClient{}
		transportMock.On(
			"DoRequest",
			ctx,
			path.Join("cbr", "currencies"),
			url.Values{},
			mock.AnythingOfType("*bytes.Buffer"),
		).Return(models.NewHTTPResponse(http.StatusOK, []byte(`invalid json`)), nil).Once()

		client := NewCbrClient(logger, transportMock)
		_, err := client.GetAllCurrencies(ctx, date)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal response")

		transportMock.AssertExpectations(t)
	})
}
