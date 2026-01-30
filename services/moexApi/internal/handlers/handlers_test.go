package handlers

import (
	"errors"
	"main/internal/service"
	"main/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	http_server = "localhost:8081"
)

var happyPathValues = service.Values{
	ShortName: service.NullString{
		Value:  "СибСтекП04",
		IsSet:  true,
		IsNull: false,
	},
	TradeDate: service.NullString{
		Value:  "2025-10-21",
		IsSet:  true,
		IsNull: false,
	},
	MaturityDate: service.NullString{
		Value:  "2027-09-28",
		IsSet:  true,
		IsNull: false,
	},
	OfferDate: service.NullString{
		Value:  "2026-04-13",
		IsSet:  true,
		IsNull: false,
	},
	BuybackDate: service.NullString{
		Value:  "2026-04-13",
		IsSet:  true,
		IsNull: false,
	},
	YieldToMaturity: service.NullFloat64{
		Value:  15.61,
		IsSet:  true,
		IsNull: false,
	},
	YieldToOffer: service.NullFloat64{
		Value:  14.3649,
		IsSet:  true,
		IsNull: false,
	},
	FaceValue: service.NullFloat64{
		Value:  1000,
		IsSet:  true,
		IsNull: false,
	},
	FaceUnit: service.NullString{
		Value:  "RUB",
		IsSet:  true,
		IsNull: false,
	},
	Duration: service.NullFloat64{
		Value:  162,
		IsSet:  true,
		IsNull: false,
	},
}

var jsonHappyPath = `{
  "SHORTNAME": {
    "value": "СибСтекП04",
    "isSet": true,
    "isNull": false
  },
  "TRADEDATE": {
    "value": "2025-10-21",
    "isSet": true,
    "isNull": false
  },
  "MATDATE": {
    "value": "2027-09-28",
    "isSet": true,
    "isNull": false
  },
  "OFFERDATE": {
    "value": "2026-04-13",
    "isSet": true,
    "isNull": false
  },
  "BUYBACKDATE": {
    "value": "2026-04-13",
    "isSet": true,
    "isNull": false
  },
  "YIELDCLOSE": {
    "value": 15.61,
    "isSet": true,
    "isNull": false
  },
  "YIELDTOOFFER": {
    "value": 14.3649,
    "isSet": true,
    "isNull": false
  },
  "FACEVALUE": {
    "value": 1000,
    "isSet": true,
    "isNull": false
  },
  "FACEUNIT": {
    "value": "RUB",
    "isSet": true,
    "isNull": false
  },
  "DURATION": {
    "value": 162,
    "isSet": true,
    "isNull": false
  }
}`

func TestGetSpecifications_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   http_server,
	}

	e := httpexpect.Default(t, u.String())

	obj := e.POST("/specifications").
		WithJSON(service.SpecificationsRequest{
			Ticker: "RU000A107209",
			Date:   time.Date(2025, time.October, 21, 0, 0, 0, 0, time.UTC),
		}).
		Expect().
		Status(200).
		JSON().Object()

	obj.Path("$.SHORTNAME.value").String().IsEqual(happyPathValues.ShortName.Value)
	obj.Path("$.TRADEDATE.value").String().IsEqual(happyPathValues.TradeDate.Value)
	obj.Path("$.MATDATE.value").String().IsEqual(happyPathValues.MaturityDate.Value)
	obj.Path("$.OFFERDATE.value").String().IsEqual(happyPathValues.OfferDate.Value)
	obj.Path("$.BUYBACKDATE.value").String().IsEqual(happyPathValues.BuybackDate.Value)
	obj.Path("$.YIELDCLOSE.value").Number().IsEqual(happyPathValues.YieldToMaturity.Value)
	obj.Path("$.YIELDTOOFFER.value").Number().IsEqual(happyPathValues.YieldToOffer.Value)
	obj.Path("$.FACEVALUE.value").Number().IsEqual(happyPathValues.FaceValue.Value)
	obj.Path("$.FACEUNIT.value").String().IsEqual(happyPathValues.FaceUnit.Value)
	obj.Path("$.DURATION.value").Number().IsEqual(happyPathValues.Duration.Value)
}

func TestGetSpecification(t *testing.T) {
	cases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		mockErr        error
	}{
		{
			name:           "successful request",
			requestBody:    `{"ticker": "RU000A107209", "date": "2025-10-21T00:00:00.00Z"}`,
			expectedStatus: 200,
			expectedBody:   jsonHappyPath,
			mockErr:        nil,
		},
		{
			name:           "service.Getspecifications err",
			requestBody:    `{"ticker": "RU000tryA107209", "date": "2025-10-21T00:00:00.00Z"}`,
			expectedStatus: 500,
			expectedBody:   `{"error": "Could not get specifications"}`,
			mockErr:        errors.New("could not find data in MOEX"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			serviceMock := mocks.NewService(t)
			if tc.name == "successful request" {
				serviceMock.On("GetSpecifications", service.SpecificationsRequest{
					Ticker: "RU000A107209",
					Date:   time.Date(2025, time.October, 21, 0, 0, 0, 0, time.UTC)}).
					Once().
					Return(happyPathValues, tc.mockErr)
			}
			if tc.name == "service.Getspecifications err" {
				serviceMock.On("GetSpecifications", service.SpecificationsRequest{
					Ticker: "RU000tryA107209",
					Date:   time.Date(2025, time.October, 21, 0, 0, 0, 0, time.UTC)}).
					Once().
					Return(service.Values{}, tc.mockErr)
			}

			handl := NewHandlers(serviceMock)

			c, rec := createTestContext(http.MethodPost, "/specifications", tc.requestBody)

			err := handl.GetSpecifications(c)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
		})
	}

}

func TestGetSpecification_InvalidJson(t *testing.T) {
	serviceMock := mocks.NewService(t)

	handler := NewHandlers(serviceMock)

	c, rec := createTestContext(http.MethodPost, "/specifications", `{invalid json`)

	err := handler.GetSpecifications(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"error":"Invalid request"}`, rec.Body.String())
	serviceMock.AssertNotCalled(t, "GetSpecifications", mock.Anything)
}

func createTestContext(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
