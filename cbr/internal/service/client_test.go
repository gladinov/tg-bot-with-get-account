package service

import (
	timezone "cbr/lib/timeZone"
	"errors"
	"fmt"
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
			transport := NewTransport(cbrHost)
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
	cases := []struct {
		name    string
		data    []byte
		want    ValCurs
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
			want:    ValCurs{},
			wantErr: errors.New("op: service.parseCurrencies, could not decode Xml file"),
		},
	}
	for _, tc := range cases {
		fmt.Println(string(xmlDataInBytes[:80]))
		transport := NewTransport(cbrHost)
		client := NewClient(transport)
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

	client := NewClient(nil)
	got, err := client.parseCurrencies(xmlData)

	require.NoError(t, err)
	require.Equal(t, "Австралийский доллар", got.Valute[0].Name)
}
