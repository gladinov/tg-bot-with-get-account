package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNormalizeDate(t *testing.T) {
	location, _ := GetMoscowLocation()
	now := time.Now()
	startDate := GetStartSingleExchangeRateRubble(location)
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
			got := NormalizeDate(tc.date, now, startDate)
			require.Equal(t, tc.want, got)
		})
	}
}
