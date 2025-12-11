package service

import (
	timezone "cbr/lib/timeZone"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	cbrHost = "www.cbr.ru"
)

func TestDoRequest(t *testing.T) {
	logg := slog.Default()
	location, _ := timezone.GetMoscowLocation()
	now := time.Now().In(location).Format(layout)
	cases := []struct {
		name    string
		date    string
		path    string
		params  url.Values
		want    []byte
		wantErr error
	}{
		{
			name:    "Err : doRequest err",
			date:    now,
			path:    path.Join("scripts", "XML_daily.asp"),
			params:  url.Values{"date_req": []string{now}},
			wantErr: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			transport := NewTransport(cbrHost, logg)
			got, err := transport.DoRequest(tc.path, tc.params)
			if tc.wantErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Contains(t, string(got), "<ValCurs")
		})
	}
}

func TestParseCurrencies(t *testing.T) {
	logg := slog.Default()
	cases := []struct {
		name    string
		data    []byte
		want    CurrenciesResponce
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
			want:    CurrenciesResponce{},
			wantErr: errors.New("op: service.parseCurrencies, could not decode Xml file"),
		},
	}
	for _, tc := range cases {
		fmt.Println(string(xmlDataInBytes[:80]))
		transport := NewTransport(cbrHost, logg)
		client := NewClient(transport, logg)
		got, err := client.parseCurrencies(tc.data)

		if tc.wantErr != nil {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
}

func TestParseCurrencies_UTF8CharsetReader(t *testing.T) {
	logg := slog.Default()
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

	client := NewClient(nil, logg)
	got, err := client.parseCurrencies(xmlData)

	require.NoError(t, err)
	require.Equal(t, "Австралийский доллар", got.Currencies[0].Name)
}

func TestNormalizeDate(t *testing.T) {
	location, _ := timezone.GetMoscowLocation()
	now := time.Now()
	startDate := timezone.GetStartSingleExchangeRateRubble(location)
	cases := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "Sucsess",
			date: now,
			want: now.Format(layout),
		},
		{
			name: "FutureDate",
			date: now.AddDate(100, 0, 0),
			want: now.Format(layout),
		},
		{
			name: "PastDate",
			date: now.AddDate(-100, 0, 0),
			want: startDate.Format(layout),
		},
		{
			name: "Date After Start Single Exchange Rate Rubble",
			date: startDate.AddDate(0, 0, 1),
			want: startDate.AddDate(0, 0, 1).Format(layout),
		},
		{
			name: "Border case",
			date: startDate,
			want: startDate.Format(layout),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeDate(tc.date, now, startDate)
			require.Equal(t, tc.want, got)
		})
	}

}
