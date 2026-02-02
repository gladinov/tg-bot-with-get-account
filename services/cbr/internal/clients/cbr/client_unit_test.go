//go:build unit

package cbr

import (
	"cbr/internal/clients/cbr/mocks"
	"cbr/internal/utils"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	cbrHost = "www.cbr.ru"
	layout  = "02/01/2006"
)

func TestGetAllCurrencies(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))

	location, _ := utils.GetMoscowLocation()
	now := time.Now().In(location)
	startDate := utils.GetStartSingleExchangeRateRubble(location)
	date := time.Now().AddDate(0, 0, -1)
	formatDate := utils.NormalizeDate(date, now, startDate)

	t.Run("sucsess", func(t *testing.T) {
		transportMock := mocks.NewHTTPTransport(t)
		transportMock.On(
			"DoRequest",
			ctx,
			"scripts/XML_daily.asp",
			mock.MatchedBy(func(v url.Values) bool {
				return v.Get("date_req") == formatDate
			}),
		).Return(xmlDataInBytes, nil).Once()
		client := NewClient(logg, transportMock)
		currResp, err := client.GetAllCurrencies(ctx, formatDate)
		require.NoError(t, err)
		require.Equal(t, xmlData, currResp)
		transportMock.AssertExpectations(t)
	})
	t.Run("DoRequest failed", func(t *testing.T) {
		transportMock := mocks.NewHTTPTransport(t)
		transportMock.On("DoRequest",
			ctx,
			"scripts/XML_daily.asp",
			mock.MatchedBy(func(v url.Values) bool {
				return v.Get("date_req") == formatDate
			}),
		).Return(nil, errors.New("could not do request")).Once()

		client := NewClient(logg, transportMock)
		currResp, err := client.GetAllCurrencies(ctx, formatDate)
		require.Error(t, err)
		require.Empty(t, currResp)
		require.ErrorContains(t, err, "could not do request")
		transportMock.AssertExpectations(t)
	})
	t.Run("parseCurrencies failed", func(t *testing.T) {
		transportMock := mocks.NewHTTPTransport(t)
		transportMock.On("DoRequest",
			ctx,
			"scripts/XML_daily.asp",
			mock.MatchedBy(func(v url.Values) bool {
				return v.Get("date_req") == formatDate
			}),
		).Return(xmlDataInBytesErr, nil).Once()

		client := NewClient(logg, transportMock)
		currResp, err := client.GetAllCurrencies(ctx, formatDate)
		require.Error(t, err)
		require.Empty(t, currResp)
		require.ErrorContains(t, err, "could not parse currencies")
		transportMock.AssertExpectations(t)
	})
}
