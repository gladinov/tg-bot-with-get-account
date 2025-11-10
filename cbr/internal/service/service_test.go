package service_test

import (
	"cbr/internal/service"
	"cbr/internal/service/mocks"
	timezone "cbr/lib/timeZone"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	layout  = "02/01/2006"
	cbrHost = "www.cbr.ru"
)

func TestGetAllCurrencies_Date(t *testing.T) {
	location, _ := timezone.GetMoscowLocation()
	startDate := timezone.GetStartSingleExchangeRateRubble(location)
	startDateStr := startDate.Format(layout)
	now := time.Now().In(location).Format(layout)
	cases := []struct {
		name     string
		date     time.Time
		wantDate string
		// wantValCurs ValCurs
		wantErr bool
	}{
		{
			name: "HappyPath",
			date: time.Date(2025, time.November, 6, 0, 0, 0, 0, time.UTC),
			// wantValCurs: happyPathCurrencies,
			wantDate: "06.11.2025",
			wantErr:  false,
		},
		{
			name: "Future date",
			date: time.Date(2025, time.November, 6, 0, 0, 0, 0, time.UTC).AddDate(100, 0, 0),
			// wantValCurs: ValCurs{
			// 	Date:   "07.11.2025",
			// 	Valute: happyPathCurrencies.Valute,
			// },
			wantDate: now,
			wantErr:  false,
		},
		{
			name:     "Past date",
			date:     time.Date(2025, time.November, 6, 0, 0, 0, 0, time.UTC).AddDate(-100, 0, 0),
			wantDate: startDateStr,
			wantErr:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			transport := service.NewTransport(cbrHost)
			client := service.NewClient(transport)

			cbr := service.NewService(client)
			_, err := cbr.GetAllCurrencies(tc.date)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			// require.Equal(t, tc.want, got)

		})
	}
}

func TestGetAllCurrencies(t *testing.T) {
	location, _ := timezone.GetMoscowLocation()
	now := time.Now().In(location)
	cases := []struct {
		name          string
		date          time.Time
		path          string
		params        url.Values
		want          service.ValCurs
		setupMock     func(*mocks.HTTPTransport, string, url.Values)
		wantErr       error
		assertNoCalls bool
	}{
		{
			name:   "Err : doRequest err",
			date:   now,
			path:   "",
			params: url.Values{},
			want:   service.ValCurs{},

			setupMock: func(mockHTTPTransport *mocks.HTTPTransport, path string, params url.Values) {
				mockHTTPTransport.On("DoRequest", mock.Anything, mock.Anything).
					Once().
					Return(nil, errors.New("op: service.DoRequest, error: could not create http.NewRequest"))
			},
			wantErr:       errors.New("op: service.GetAllCurrencies, error: could not do request"),
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockHTTPTransport := mocks.NewHTTPTransport(t)
			tc.setupMock(mockHTTPTransport, tc.path, tc.params)
			client := service.NewClient(mockHTTPTransport)
			cbr := service.NewService(client)
			got, err := cbr.GetAllCurrencies(tc.date)
			if tc.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.want, got)
				require.Equal(t, tc.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantErr, err)
			if tc.assertNoCalls {
				mockHTTPTransport.AssertNotCalled(t, "DoRequest")
			} else {
				mockHTTPTransport.AssertExpectations(t)
			}
		})
	}
}
