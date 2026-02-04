//go:build unit

package cbr

import (
	"cbr/internal/models"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
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

		got, err := parseCurrencies(ctx, logg, tc.data)

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

	got, err := parseCurrencies(ctx, logg, xmlData)

	require.NoError(t, err)
	require.Equal(t, "Австралийский доллар", got.Currencies[0].Name)
}
