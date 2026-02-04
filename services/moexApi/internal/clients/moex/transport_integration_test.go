//go:build integration

package moex

import (
	"context"
	"io"
	"log/slog"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	moexHost = "iss.moex.com"
)

func TestDoRequest(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))
	cases := []struct {
		name        string
		path        string
		query       url.Values
		host        string
		expected    []byte
		expectedErr bool
	}{
		{
			name:        "Correct RU000A107209",
			path:        path.Join("iss", "history", "engines", "stock", "markets", "bonds", "sessions", "3", "securities", "RU000A107209"+".json"),
			query:       createParams(time.Now()),
			host:        moexHost,
			expectedErr: false,
		},
		{
			name:        "Err client.Do",
			path:        path.Join("iss", "history", "engines", "stock", "markets", "bonds", "sessions", "3", "securities", "RU000A107209"+".json"),
			query:       createParams(time.Now()),
			host:        "oitgjreoji0043o3ilkdfng",
			expectedErr: true,
		},
		{
			name:        "Err http.NewRequest",
			path:        path.Join("iss", "history", "engines", "stock", "markets", "bonds", "sessions", "3", "securities", "RU000A107209"+".json"),
			query:       createParams(time.Now()),
			host:        "oitgjreoji0043o3ilkdfng\b",
			expectedErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			transport := NewTransport(logg, tc.host)
			_, err := transport.DoRequest(ctx, tc.path, tc.query)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func createParams(date time.Time) url.Values {
	formatDate := date.Format(layout)
	params := url.Values{}
	params.Add("limit", "1")
	params.Add("iss.meta", "off")
	params.Add("history.columns", "TRADEDATE,MATDATE,OFFERDATE,BUYBACKDATE,YIELDCLOSE,YIELDTOOFFER,FACEVALUE,FACEUNIT,DURATION, SHORTNAME")
	params.Add("limit", "1")
	params.Add("from", formatDate)
	params.Add("to", formatDate)
	return params
}
