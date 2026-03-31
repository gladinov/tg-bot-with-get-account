//go:build integration

package cbr

import (
	"cbr/internal/utils"
	"context"
	"io"
	"log/slog"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoRequest(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))
	location, _ := utils.GetMoscowLocation()
	now := time.Now().In(location).Format(layout)
	cases := []struct {
		name        string
		host        string
		date        string
		path        string
		params      url.Values
		want        []byte
		wantErr     bool
		errContains string
	}{
		{
			name:    "Success",
			host:    cbrHost,
			date:    now,
			path:    path.Join("scripts", "XML_daily.asp"),
			params:  url.Values{"date_req": []string{now}},
			wantErr: false,
		},
		{
			name:        "Fail: invalid host",
			host:        "invalid host",
			path:        "scripts/XML_daily.asp",
			params:      url.Values{},
			wantErr:     true,
			errContains: "could not create http.NewRequest",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			transport := NewTransport(logg, tc.host)
			got, err := transport.DoRequest(ctx, tc.path, tc.params)
			if tc.wantErr {
				require.ErrorContains(t, err, tc.errContains)
				return
			}
			require.NoError(t, err)
			require.Contains(t, string(got), "<ValCurs")
		})
	}
}
