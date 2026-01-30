//go:build unit

package cbr

import (
	"cbr/internal/clients/cbr/mocks"
	"cbr/internal/models"
	"cbr/internal/utils"
	"context"
	"errors"
	"fmt"
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

func TestParseCurrencies(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))
	cases := []struct {
		name    string
		data    []byte
		want    models.CurrenciesResponce
		wantErr error
	}{
		{
			name:    "Sucsess",
			data:    xmlDataInBytes,
			want:    xmlData,
			wantErr: nil,
		},
		{
			name:    "Err: Incorrect data in byte",
			data:    xmlDataInBytesErr,
			want:    models.CurrenciesResponce{},
			wantErr: errors.New("op: service.parseCurrencies, could not decode Xml file"),
		},
	}
	for _, tc := range cases {
		fmt.Println(string(xmlDataInBytes[:80]))
		transport := NewTransport(logg, cbrHost)
		client := NewClient(logg, transport)
		got, err := client.parseCurrencies(ctx, tc.data)

		if tc.wantErr != nil {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
}

func TestParseCurrencies_UTF8CharsetReader(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))
	xmlData := []byte(`<?xml version="1.0" encoding="koi8-r"?>
<ValCurs Date="03.03.1995" name="Foreign Currency Market">
    <Valute ID="R01010">
        <NumCode>036</NumCode>
        <CharCode>AUD</CharCode>
        <Nominal>1</Nominal>
        <Name>Австралийский доллар</Name>
        <Value>3334,8200</Value>
        <VunitRate>3334,82</VunitRate>
    </Valute>
</ValCurs>`)

	client := NewClient(logg, nil)
	got, err := client.parseCurrencies(ctx, xmlData)

	require.NoError(t, err)
	require.Equal(t, "Австралийский доллар", got.Currencies[0].Name)
}

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
