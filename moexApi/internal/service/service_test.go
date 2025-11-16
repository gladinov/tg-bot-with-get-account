package service

import (
	"net/http"
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
			service := NewSpecificationService(tc.host)

			_, err := service.DoRequest(tc.path, tc.query)
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

func TestGetSpecifications(t *testing.T) {
	cases := []struct {
		name    string
		input   SpecificationsRequest
		want    Values
		wantErr bool
	}{
		{
			name: "succses",
			input: SpecificationsRequest{
				Ticker: "RU000A107209",
				Date:   time.Date(2025, time.October, 21, 0, 0, 0, 0, time.UTC),
			},
			want: Values{
				ShortName: NullString{
					Value:  "СибСтекП04",
					IsSet:  true,
					IsNull: false,
				},
				TradeDate: NullString{
					Value:  "2025-10-21",
					IsSet:  true,
					IsNull: false,
				},
				MaturityDate: NullString{
					Value:  "2027-09-28",
					IsSet:  true,
					IsNull: false,
				},
				OfferDate: NullString{
					Value:  "2026-04-13",
					IsSet:  true,
					IsNull: false,
				},
				BuybackDate: NullString{
					Value:  "2026-04-13",
					IsSet:  true,
					IsNull: false,
				},
				YieldToMaturity: NullFloat64{
					Value:  15.61,
					IsSet:  true,
					IsNull: false,
				},
				YieldToOffer: NullFloat64{
					Value:  14.3649,
					IsSet:  true,
					IsNull: false,
				},
				FaceValue: NullFloat64{
					Value:  1000,
					IsSet:  true,
					IsNull: false,
				},
				FaceUnit: NullString{
					Value:  "RUB",
					IsSet:  true,
					IsNull: false,
				},
				Duration: NullFloat64{
					Value:  162,
					IsSet:  true,
					IsNull: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Incorrect ticker",
			input: SpecificationsRequest{
				Ticker: "oigheorgih",
				Date:   time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := SpecificationService{
				host:   moexHost,
				client: http.Client{},
			}
			got, err := service.GetSpecifications(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertNullStringEqual(t, "ShortName", tc.want.ShortName, got.ShortName)
			assertNullStringEqual(t, "TradeDate", tc.want.TradeDate, got.TradeDate)
			assertNullStringEqual(t, "MaturityDate", tc.want.MaturityDate, got.MaturityDate)
			assertNullStringEqual(t, "OfferDate", tc.want.OfferDate, got.OfferDate)
			assertNullStringEqual(t, "BuybackDate", tc.want.BuybackDate, got.BuybackDate)
			assertNullFloat64Equal(t, "YieldToMaturity", tc.want.YieldToMaturity, got.YieldToMaturity)
			assertNullFloat64Equal(t, "YieldToOffer", tc.want.YieldToOffer, got.YieldToOffer)
			assertNullFloat64Equal(t, "FaceValue", tc.want.FaceValue, got.FaceValue)
			assertNullStringEqual(t, "FaceUnit", tc.want.FaceUnit, got.FaceUnit)
			assertNullFloat64Equal(t, "Duration", tc.want.Duration, got.Duration)
		})

	}
}

func TestGetSpecifications_Date(t *testing.T) {
	cases := []struct {
		name    string
		input   SpecificationsRequest
		want    Values
		wantErr bool
	}{
		{
			name: "Data from future",
			input: SpecificationsRequest{
				Ticker: "RU000A107209",
				Date:   time.Now().AddDate(25, 0, 0),
			},
			want: Values{

				TradeDate: NullString{
					Value:  time.Now().Format(layout),
					IsSet:  true,
					IsNull: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Data from future",
			input: SpecificationsRequest{
				Ticker: "RU000A107209",
				Date:   time.Now().AddDate(-25, 0, 0),
			},
			want: Values{},

			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := SpecificationService{
				host:   moexHost,
				client: http.Client{},
			}
			_, err := service.GetSpecifications(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})

	}
}
