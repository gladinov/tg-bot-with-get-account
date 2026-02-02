//go:build unit

package moex

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"main/internal/clients/moex/mocks"
	"main/internal/models"
	"main/internal/testdata/factories"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetSpecification(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))
	date := time.Date(2025, 11, 16, 0, 0, 0, 0, time.UTC)
	ticker := "OFZ26238"

	t.Run("Success", func(t *testing.T) {
		transportMock := mocks.NewTransportClient(t)

		body := factories.NewSpecificationsResponseJSON()
		want := factories.NewSpecificationsResponse()

		transportMock.
			On("DoRequest", ctx,
				mock.MatchedBy(func(p string) bool {
					return strings.HasSuffix(p, ticker+".json")
				}),
				mock.MatchedBy(func(v url.Values) bool {
					return v.Get("from") == date.Format(layout) &&
						v.Get("to") == date.Format(layout)
				}),
			).Return(body, nil).Once()
		moexClient := NewMoexClient(logg, transportMock)
		got, err := moexClient.GetSpecifications(ctx, ticker, date)
		require.NoError(t, err)
		require.Equal(t, want, got)

		transportMock.AssertExpectations(t)
	})
	t.Run("DoRequest err", func(t *testing.T) {
		transportMock := mocks.NewTransportClient(t)

		errContains := "failed DoRequest"

		transportMock.On("DoRequest", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("url.Values")).
			Return(nil, errors.New("could not do request")).Once()
		moexClient := NewMoexClient(logg, transportMock)
		_, err := moexClient.GetSpecifications(ctx, "ticker", date)
		require.Error(t, err)
		require.ErrorContains(t, err, errContains)

		transportMock.AssertExpectations(t)
	})
	t.Run("Invalid JSON", func(t *testing.T) {
		transportMock := mocks.NewTransportClient(t)

		transportMock.On("DoRequest", ctx, mock.Anything, mock.Anything).
			Return([]byte(`invalid json`), nil).Once()

		client := NewMoexClient(logg, transportMock)
		_, err := client.GetSpecifications(ctx, ticker, date)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed unmarshall json")

		transportMock.AssertExpectations(t)
	})

	t.Run("Empty History", func(t *testing.T) {
		transportMock := mocks.NewTransportClient(t)

		emptyJSON := []byte(`{"history":{"data":[]}}`)
		want := factories.NewSpecificationsResponse()
		want.History.Data = []models.Values{} // пустой массив

		transportMock.On("DoRequest", ctx, mock.Anything, mock.Anything).
			Return(emptyJSON, nil).Once()

		client := NewMoexClient(logg, transportMock)
		got, err := client.GetSpecifications(ctx, ticker, date)
		require.NoError(t, err)
		require.Equal(t, want, got)

		transportMock.AssertExpectations(t)
	})
}
