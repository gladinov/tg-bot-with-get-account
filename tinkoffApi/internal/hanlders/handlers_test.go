package hanlders

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"tinkoffApi/internal/service"
	"tinkoffApi/internal/service/mocks"
	testhelpfunc "tinkoffApi/lib/testHelpFunc"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"main.go/pkg/app"
)

const (
	http_server = "localhost:8082"
)

func TestAuth(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name        string
		token       string
		tokenPrefix string
		wantHeader  bool
		want        string
		wantErr     error
	}{
		{
			name:        "HappyPath",
			token:       tokens.OnlyReadToken,
			tokenPrefix: "Bearer",
			wantHeader:  true,
			want:        tokens.OnlyReadToken,
			wantErr:     nil,
		},
		{
			name:        "Err: Incorrect prefix",
			token:       tokens.OnlyReadToken,
			tokenPrefix: "Bearetgioerjoi",
			wantHeader:  true,
			want:        "",
			wantErr:     errInvalidAuthFormat,
		},
		{
			name:        "Err: Zero token and zero prefix",
			token:       "",
			tokenPrefix: "",
			wantHeader:  true,
			want:        "",
			wantErr:     errHeaderRequierd,
		},
		{
			name:        "Err: Don't install auth header",
			token:       "",
			tokenPrefix: "",
			wantHeader:  false,
			want:        "",
			wantErr:     errHeaderRequierd,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := createTestContextFoAuth(tc.token, http.MethodGet, "/accounts", "", tc.tokenPrefix, tc.wantHeader)
			got, err := auth(c)
			if tc.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.want, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetAccounts(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name           string
		token          string
		expectedStatus int
		expectedBody   string
		setupMocks     func(*mocks.PortfolioService, *mocks.AnalyticsService, *mocks.InstrumentService)
	}{
		{
			name:           "HappyPath(token:OnlyRead)",
			token:          tokens.OnlyReadToken,
			expectedStatus: http.StatusOK,
			setupMocks: func(ps *mocks.PortfolioService, as *mocks.AnalyticsService, is *mocks.InstrumentService) {
				ps.On("GetClient", mock.Anything, mock.AnythingOfType("string")).
					Once().
					Return(nil)

				ps.On("GetAccounts").
					Once().
					Return(accountsServiceHappyPath, nil)
			},
			expectedBody: accountHappyPath,
		},
		{
			name:           "Err: Incorerct token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "empty token"}`,
		},
		{
			name:           "Err: GetClient err",
			token:          tokens.OnlyReadToken,
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(ps *mocks.PortfolioService, as *mocks.AnalyticsService, is *mocks.InstrumentService) {
				ps.On("GetClient", mock.Anything, mock.AnythingOfType("string")).
					Once().
					Return(errors.New("op:sevrice.AnalyticsServiceClient.GetClient, err: can't connect with tinkoffApi client"))

			},
			expectedBody: `{"error": "incorrect token"}`,
		},
		{
			name:           "Err: GetAccount err",
			token:          tokens.OnlyReadToken,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(ps *mocks.PortfolioService, as *mocks.AnalyticsService, is *mocks.InstrumentService) {
				ps.On("GetClient", mock.Anything, mock.AnythingOfType("string")).
					Once().
					Return(nil)

				ps.On("GetAccounts").
					Once().
					Return(emptyAccounts, errors.New("op:sevrice.GetAccounts, err: could not get accounts"))
			},
			expectedBody: `{"error": "Could not get accounts"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			portfolioServiceMock := mocks.NewPortfolioService(t)
			analyticServiceMock := mocks.NewAnalyticsService(t)
			instrumentServiceMock := mocks.NewInstrumentService(t)

			if tc.setupMocks != nil {
				tc.setupMocks(portfolioServiceMock, analyticServiceMock, instrumentServiceMock)
			}
			serviceClient := service.NewService(analyticServiceMock, portfolioServiceMock, instrumentServiceMock)
			handler := NewHandlers(serviceClient)
			c, rec := createTestContextWithoutBody(tc.token, http.MethodGet, "/accounts")

			err := handler.GetAccounts(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			portfolioServiceMock.AssertExpectations(t)
		})
	}

}

func TestGetPortfolio(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name           string
		token          string
		request        any
		expectedStatus int
		setupMocks     func(*mocks.PortfolioService, string, string, int64)
		expectedBody   string
		assertNoCalls  bool
	}{
		{
			name:           "HappyPath",
			token:          tokens.OnlyReadToken,
			request:        portfolioRequest,
			expectedStatus: http.StatusOK,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, accountStatus int64) {
				ps.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				ps.On("GetPortfolio", service.PortfolioRequest{AccountID: accountId, AccountStatus: accountStatus}).
					Once().
					Return(TestPortfolio, nil)
			},
			expectedBody:  portffolioHappyPath,
			assertNoCalls: false,
		},
		{
			name:           "Err: Incorect auth",
			token:          "",
			request:        portfolioRequest,
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, accountStatus int64) {
			},
			expectedBody:  `{"error": "incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:           "Err: Incorrect request body",
			token:          tokens.OnlyReadToken,
			request:        `{"invalid json"}`,
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, accountStatus int64) {
			},
			expectedBody:  `{"error": "invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:           "Err: GetClient err",
			token:          tokens.OnlyReadToken,
			request:        portfolioRequest,
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, accountStatus int64) {
				ps.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.PortfolioServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:           "Err: GetPortfolio err ",
			token:          tokens.OnlyReadToken,
			request:        portfolioRequest,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, accountStatus int64) {
				ps.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				ps.On("GetPortfolio", service.PortfolioRequest{AccountID: accountId, AccountStatus: accountStatus}).
					Once().
					Return(service.Portfolio{}, errors.New("op: service.GetPortfolio, error: can't get portifolio positions from tinkoff Api"))
			},
			expectedBody:  `{"error": "Could not get portfolio"}`,
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			portfolioServiceMock := mocks.NewPortfolioService(t)
			ananlyticServiceMock := mocks.NewAnalyticsService(t)
			instrumentServiceMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMocks(portfolioServiceMock, tc.token, "", 0)
			case []byte:
				jsonData = body
				tc.setupMocks(portfolioServiceMock, tc.token, "", 0)
			case service.PortfolioRequest:
				jsonData, _ = json.Marshal(body)
				tc.setupMocks(portfolioServiceMock, tc.token, body.AccountID, body.AccountStatus)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMocks(portfolioServiceMock, tc.token, "", 0)
			}

			serviceClient := service.NewService(ananlyticServiceMock, portfolioServiceMock, instrumentServiceMock)
			handler := NewHandlers(serviceClient)

			c, rec := createTestContext(tc.token, http.MethodPost, "/portfolio", jsonData)

			err := handler.GetPortfolio(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				portfolioServiceMock.AssertNotCalled(t, "GetClient")
				portfolioServiceMock.AssertNotCalled(t, "GetPortfolio")
			} else {
				portfolioServiceMock.AssertExpectations(t)
			}
		})
	}

}

func TestGetOperations(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name           string
		token          string
		request        any
		expectedStatus int
		setupMocks     func(*mocks.PortfolioService, string, string, time.Time)
		expectedBody   string
		assertNoCalls  bool
	}{
		{
			name:           "HappyPath",
			token:          tokens.OnlyReadToken,
			request:        operationRequest,
			expectedStatus: http.StatusOK,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, date time.Time) {
				ps.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				ps.On("MakeSafeGetOperationsRequest", service.OperationsRequest{AccountID: accountId, Date: date}).
					Once().
					Return(happyPathOperationResponse, nil)
			},
			expectedBody:  happyPathOperationResponseInBytes,
			assertNoCalls: false,
		},
		{
			name:           "Err: Incorect auth",
			token:          "",
			request:        operationRequest,
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, date time.Time) {
			},
			expectedBody:  `{"error": "incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:           "Err: Incorrect request body",
			token:          tokens.OnlyReadToken,
			request:        `{"invalid json"}`,
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, date time.Time) {
			},
			expectedBody:  `{"error": "invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:           "Err: GetClient err",
			token:          tokens.OnlyReadToken,
			request:        operationRequest,
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, date time.Time) {
				ps.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.PortfolioServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:           "Err: GetOperation err ",
			token:          tokens.OnlyReadToken,
			request:        operationBadRequest,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(ps *mocks.PortfolioService, token string, accountId string, date time.Time) {
				ps.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				ps.On("MakeSafeGetOperationsRequest", service.OperationsRequest{AccountID: accountId, Date: date}).
					Once().
					Return(nil, errors.New("op: tinkoffApi.GetOperations, error: empty account ID"))
			},
			expectedBody:  `{"error": "Could not get operations"}`,
			assertNoCalls: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			portfolioServiceMock := mocks.NewPortfolioService(t)
			analyticServiceMock := mocks.NewAnalyticsService(t)
			instrumentServiceMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMocks(portfolioServiceMock, tc.token, "", time.Now())
			case []byte:
				jsonData = body
				tc.setupMocks(portfolioServiceMock, tc.token, "", time.Now())
			case service.OperationsRequest:
				jsonData, _ = json.Marshal(body)
				tc.setupMocks(portfolioServiceMock, tc.token, body.AccountID, body.Date)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMocks(portfolioServiceMock, tc.token, "", time.Now())
			}

			serviceClient := service.NewService(analyticServiceMock, portfolioServiceMock, instrumentServiceMock)
			handler := NewHandlers(serviceClient)

			c, rec := createTestContext(tc.token, http.MethodPost, "/operations", jsonData)

			err := handler.GetOperations(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				portfolioServiceMock.AssertNotCalled(t, "GetClient")
				portfolioServiceMock.AssertNotCalled(t, "MakeSafeGetOperationsRequest")
			} else {
				portfolioServiceMock.AssertExpectations(t)
			}
		})
	}
}

func TestGetAllAssetUids(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name           string
		token          string
		expectedStatus int
		setupMocks     func(*mocks.AnalyticsService, string)
		expectedBody   string
		assertNoCalls  bool
	}{
		{
			name:  "HappyPath",
			token: tokens.OnlyReadToken,

			expectedStatus: http.StatusOK,
			setupMocks: func(as *mocks.AnalyticsService, token string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				as.On("GetAllAssetUids").
					Once().
					Return(happyPathAllAssetsUid, nil)
			},
			expectedBody:  happyPathAllAssetsUidInBytes,
			assertNoCalls: false,
		},
		{
			name:           "Err: Incorect auth",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			setupMocks:     func(as *mocks.AnalyticsService, token string) {},
			expectedBody:   `{"error": "incorrect token"}`,
			assertNoCalls:  true,
		},
		{
			name:           "Err: GetClient err",
			token:          tokens.OnlyReadToken,
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(as *mocks.AnalyticsService, token string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.AnalyticsServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:           "Err: GetAllAssetUids err ",
			token:          tokens.OnlyReadToken,
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(as *mocks.AnalyticsService, token string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				as.On("GetAllAssetUids").
					Once().
					Return(nil, errors.New("op: service.GetAllAssetUids, error: could not get assets uid"))
			},
			expectedBody:  `{"error": "Could not get assets uids"}`,
			assertNoCalls: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			portfolioServiceMock := mocks.NewPortfolioService(t)
			analyticServiceMock := mocks.NewAnalyticsService(t)
			instrumentServiceMock := mocks.NewInstrumentService(t)
			tc.setupMocks(analyticServiceMock, tc.token)
			serviceClient := service.NewService(analyticServiceMock, portfolioServiceMock, instrumentServiceMock)
			handler := NewHandlers(serviceClient)

			c, rec := createTestContextWithoutBody(tc.token, http.MethodGet, "/allassetsuid")

			err := handler.GetAllAssetUids(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				portfolioServiceMock.AssertNotCalled(t, "GetClient")
				portfolioServiceMock.AssertNotCalled(t, "GetAllAssetUids")
			} else {
				portfolioServiceMock.AssertExpectations(t)
			}
		})
	}
}

func TestGetBondsActions(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name           string
		token          string
		request        any
		expectedStatus int
		setupMocks     func(*mocks.AnalyticsService, string, string)
		expectedBody   string
		assertNoCalls  bool
	}{
		{
			name:           "HappyPath",
			token:          tokens.OnlyReadToken,
			request:        service.BondsActionsReq{InstrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			expectedStatus: http.StatusOK,
			setupMocks: func(as *mocks.AnalyticsService, token string, instrumentUid string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				as.On("GetBondsActions", instrumentUid).
					Once().
					Return(happyPathBondsActions, nil)
			},
			expectedBody:  happyPathBondsActionsInBytes,
			assertNoCalls: false,
		},
		{
			name:           "Err: Incorect auth",
			token:          "",
			request:        service.BondsActionsReq{InstrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			expectedStatus: http.StatusUnauthorized,
			setupMocks: func(as *mocks.AnalyticsService, token string, instrumentUid string) {

			},
			expectedBody:  `{"error": "incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:           "Err: Incorrect request body",
			token:          tokens.OnlyReadToken,
			request:        `{"invalid json"}`,
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(as *mocks.AnalyticsService, token string, instrumentUid string) {

			},
			expectedBody:  `{"error": "Invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:           "Err: GetClient err",
			token:          tokens.OnlyReadToken,
			request:        service.BondsActionsReq{InstrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(as *mocks.AnalyticsService, token string, instrumentUid string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.AnalyticsServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:           "Err: GetBondActions err",
			token:          tokens.OnlyReadToken,
			request:        service.BondsActionsReq{InstrumentUid: "000029ae-00a2-441c-a9bf-9bfd2988e706"},
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(as *mocks.AnalyticsService, token string, instrumentId string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				as.On("GetBondsActions", instrumentId).
					Once().
					Return(service.BondIdentIdentifiers{}, errors.New("op: service.GetBondActions, error: could not get bond actions"))
			},
			expectedBody:  `{"error": "Could not get bond identificators"}`,
			assertNoCalls: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			portfolioServiceMock := mocks.NewPortfolioService(t)
			analyticServiceMock := mocks.NewAnalyticsService(t)
			instrumentServiceMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMocks(analyticServiceMock, tc.token, "")
			case []byte:
				jsonData = body
				tc.setupMocks(analyticServiceMock, tc.token, "")
			case service.BondsActionsReq:
				jsonData, _ = json.Marshal(body)
				tc.setupMocks(analyticServiceMock, tc.token, body.InstrumentUid)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMocks(analyticServiceMock, tc.token, "")
			}

			serviceClient := service.NewService(analyticServiceMock, portfolioServiceMock, instrumentServiceMock)
			handler := NewHandlers(serviceClient)

			c, rec := createTestContext(tc.token, http.MethodPost, "/bondactions", jsonData)

			err := handler.GetBondsActions(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				portfolioServiceMock.AssertNotCalled(t, "GetClient")
				portfolioServiceMock.AssertNotCalled(t, "GetBondsActions")
			} else {
				portfolioServiceMock.AssertExpectations(t)
			}
		})
	}
}

func TestGetLastPriceInPersentageToNominal(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name          string
		token         string
		request       any
		setupMock     func(*mocks.AnalyticsService, string, string)
		expectedCode  int
		expectedBody  string
		assertNoCalls bool
	}{
		{
			name:    "HappyPath",
			token:   tokens.OnlyReadToken,
			request: service.LastPriceReq{InstrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			setupMock: func(as *mocks.AnalyticsService, token string, instrumentUid string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				as.On("GetLastPriceInPersentageToNominal", instrumentUid).
					Once().Return(happyPathLastPrice, nil)
			},
			expectedCode:  http.StatusOK,
			expectedBody:  happyPathLastPriceInBytes,
			assertNoCalls: false,
		},
		{
			name:          "Err: Auth err",
			token:         "",
			request:       service.LastPriceReq{InstrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			setupMock:     func(as *mocks.AnalyticsService, token string, instrumentUid string) {},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"error": "Incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:          "Err: Invalid request",
			token:         tokens.OnlyReadToken,
			request:       `{"invalid request"}`,
			setupMock:     func(as *mocks.AnalyticsService, token string, instrumentUid string) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "Invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:    "Err: GetClient err ",
			token:   tokens.OnlyReadToken,
			request: service.LastPriceReq{InstrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			setupMock: func(as *mocks.AnalyticsService, token string, instrumentUid string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.AnalyticsServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:    "Err: GetLastPriceInPersentageToNominal err",
			token:   tokens.OnlyReadToken,
			request: service.LastPriceReq{InstrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			setupMock: func(as *mocks.AnalyticsService, token string, instrumentUid string) {
				as.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				as.On("GetLastPriceInPersentageToNominal", instrumentUid).
					Once().
					Return(service.LastPriceResponse{}, errors.New("tinkoffApi.GetLastPriceInPercentageToNominal: instrumentUid could not be empty string"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"error": "Could not get last price"}`,
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			analyticMock := mocks.NewAnalyticsService(t)
			portfolioMock := mocks.NewPortfolioService(t)
			instrumentMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMock(analyticMock, tc.token, "")
			case []byte:
				jsonData = body
				tc.setupMock(analyticMock, tc.token, "")
			case service.LastPriceReq:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(analyticMock, tc.token, body.InstrumentUid)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(analyticMock, tc.token, "")
			}

			serviceMocks := service.NewService(analyticMock, portfolioMock, instrumentMock)
			handler := NewHandlers(serviceMocks)

			c, rec := createTestContext(tc.token, http.MethodPost, "/lastprice", jsonData)

			err := handler.GetLastPriceInPersentageToNominal(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				analyticMock.AssertNotCalled(t, "GetClient")
				analyticMock.AssertNotCalled(t, "GetLastPriceInPersentageToNominal")
			} else {
				analyticMock.AssertExpectations(t)
			}

		})
	}
}

func TestFindBy(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name          string
		token         string
		request       any
		setupMock     func(*mocks.InstrumentService, string, string)
		expectedCode  int
		expectedBody  string
		assertNoCalls bool
	}{
		{
			name:    "HappyPath",
			token:   tokens.OnlyReadToken,
			request: service.FindByReq{Query: "e80d1280-d512-4755-b48b-1187fd6cb2d8"},
			setupMock: func(is *mocks.InstrumentService, token string, query string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("FindBy", query).
					Once().Return(happpyPathFindBy, nil)
			},
			expectedCode:  http.StatusOK,
			expectedBody:  happpyPathFindByInBytes,
			assertNoCalls: false,
		},
		{
			name:          "Err: Auth err",
			token:         "",
			request:       service.FindByReq{Query: "e80d1280-d512-4755-b48b-1187fd6cb2d8"},
			setupMock:     func(is *mocks.InstrumentService, token string, query string) {},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"error": "Incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:          "Err: Invalid request",
			token:         tokens.OnlyReadToken,
			request:       `{"invalid request"}`,
			setupMock:     func(is *mocks.InstrumentService, token string, query string) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "Invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:    "Err: GetClient err ",
			token:   tokens.OnlyReadToken,
			request: service.FindByReq{Query: "e80d1280-d512-4755-b48b-1187fd6cb2d8"},
			setupMock: func(is *mocks.InstrumentService, token string, query string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.InstrumentServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:    "Err: FindBy err",
			token:   tokens.OnlyReadToken,
			request: service.FindByReq{Query: "e80d1280-d512-4755-b48b-1187fd6cb2d8"},
			setupMock: func(is *mocks.InstrumentService, token string, query string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("FindBy", query).
					Once().
					Return(nil, errors.New("op: service.FindBy, error: could not find instrument by query"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"error": "Could not get instruments"}`,
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			analyticMock := mocks.NewAnalyticsService(t)
			portfolioMock := mocks.NewPortfolioService(t)
			instrumentMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMock(instrumentMock, tc.token, "")
			case []byte:
				jsonData = body
				tc.setupMock(instrumentMock, tc.token, "")
			case service.FindByReq:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, body.Query)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, "")
			}

			serviceMocks := service.NewService(analyticMock, portfolioMock, instrumentMock)
			handler := NewHandlers(serviceMocks)

			c, rec := createTestContext(tc.token, http.MethodPost, "/findby", jsonData)

			err := handler.FindBy(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				analyticMock.AssertNotCalled(t, "GetClient")
				analyticMock.AssertNotCalled(t, "FindBy")
			} else {
				analyticMock.AssertExpectations(t)
			}

		})
	}
}

func TestGetBondByUid(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name          string
		token         string
		request       any
		setupMock     func(*mocks.InstrumentService, string, string)
		expectedCode  int
		expectedBody  string
		assertNoCalls bool
	}{
		{
			name:    "HappyPath",
			token:   tokens.OnlyReadToken,
			request: service.BondReq{Uid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			setupMock: func(is *mocks.InstrumentService, token string, uid string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetBondByUid", uid).
					Once().Return(happpyPathBondBy, nil)
			},
			expectedCode:  http.StatusOK,
			expectedBody:  happpyPathBondByInBytes,
			assertNoCalls: false,
		},
		{
			name:          "Err: Auth err",
			token:         "",
			request:       service.BondReq{Uid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			setupMock:     func(is *mocks.InstrumentService, token string, uid string) {},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"error": "Incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:          "Err: Invalid request",
			token:         tokens.OnlyReadToken,
			request:       `{"invalid request"}`,
			setupMock:     func(is *mocks.InstrumentService, token string, uid string) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "Invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:    "Err: GetClient err ",
			token:   tokens.OnlyReadToken,
			request: service.BondReq{Uid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2"},
			setupMock: func(is *mocks.InstrumentService, token string, uid string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.InstrumentServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:    "Err: GetBondByUid err",
			token:   tokens.OnlyReadToken,
			request: service.BondReq{Uid: "070d82ad-e9ewetwtwgtg0-41e4-8eca-cbe9f5830db2"},
			setupMock: func(is *mocks.InstrumentService, token string, uid string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetBondByUid", uid).
					Once().
					Return(service.Bond{}, errors.New("incorrect uid: can't be empty string"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"error": "Could not get bond"}`,
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			analyticMock := mocks.NewAnalyticsService(t)
			portfolioMock := mocks.NewPortfolioService(t)
			instrumentMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMock(instrumentMock, tc.token, "")
			case []byte:
				jsonData = body
				tc.setupMock(instrumentMock, tc.token, "")
			case service.BondReq:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, body.Uid)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, "")
			}

			serviceMocks := service.NewService(analyticMock, portfolioMock, instrumentMock)
			handler := NewHandlers(serviceMocks)

			c, rec := createTestContext(tc.token, http.MethodPost, "/bond", jsonData)

			err := handler.GetBondBy(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				analyticMock.AssertNotCalled(t, "GetClient")
				analyticMock.AssertNotCalled(t, "GetBondByUid")
			} else {
				analyticMock.AssertExpectations(t)
			}

		})
	}
}

func TestGetCurrencyBy(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name          string
		token         string
		request       any
		setupMock     func(*mocks.InstrumentService, string, string)
		expectedCode  int
		expectedBody  string
		assertNoCalls bool
	}{
		{
			name:    "HappyPath",
			token:   tokens.OnlyReadToken,
			request: service.CurrencyReq{Figi: "CNY000TODTOM"},
			setupMock: func(is *mocks.InstrumentService, token string, figi string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetCurrencyBy", figi).
					Once().Return(happpyPathCurrencyBy, nil)
			},
			expectedCode:  http.StatusOK,
			expectedBody:  happpyPathCurrencyByInBytes,
			assertNoCalls: false,
		},
		{
			name:          "Err: Auth err",
			token:         "",
			request:       service.CurrencyReq{Figi: "CNY000TODTOM"},
			setupMock:     func(is *mocks.InstrumentService, token string, uid string) {},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"error": "Incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:          "Err: Invalid request",
			token:         tokens.OnlyReadToken,
			request:       `{"invalid request"}`,
			setupMock:     func(is *mocks.InstrumentService, token string, uid string) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "Invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:    "Err: GetClient err ",
			token:   tokens.OnlyReadToken,
			request: service.CurrencyReq{Figi: "CNY000TODTOM"},
			setupMock: func(is *mocks.InstrumentService, token string, uid string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.InstrumentServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:    "Err: GetGetCurrencyBy err",
			token:   tokens.OnlyReadToken,
			request: service.CurrencyReq{Figi: "CNY000TODTOM"},
			setupMock: func(is *mocks.InstrumentService, token string, figi string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetCurrencyBy", figi).
					Once().
					Return(service.Currency{}, errors.New("op:service.GetCurrencyBy, error: can't get curency by figi"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"error": "Could not get currecny"}`,
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			analyticMock := mocks.NewAnalyticsService(t)
			portfolioMock := mocks.NewPortfolioService(t)
			instrumentMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMock(instrumentMock, tc.token, "")
			case []byte:
				jsonData = body
				tc.setupMock(instrumentMock, tc.token, "")
			case service.CurrencyReq:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, body.Figi)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, "")
			}

			serviceMocks := service.NewService(analyticMock, portfolioMock, instrumentMock)
			handler := NewHandlers(serviceMocks)

			c, rec := createTestContext(tc.token, http.MethodPost, "/currecny", jsonData)

			err := handler.GetCurrencyBy(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				analyticMock.AssertNotCalled(t, "GetClient")
				analyticMock.AssertNotCalled(t, "GetCurrencyBy")
			} else {
				analyticMock.AssertExpectations(t)
			}

		})
	}
}

func TestGetShareCurrencyBy(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name          string
		token         string
		request       any
		setupMock     func(*mocks.InstrumentService, string, string)
		expectedCode  int
		expectedBody  string
		assertNoCalls bool
	}{
		{
			name:    "HappyPath",
			token:   tokens.OnlyReadToken,
			request: service.ShareCurrencyByRequest{Figi: "BBG007N0Z367"},
			setupMock: func(is *mocks.InstrumentService, token string, figi string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetShareCurrencyBy", figi).
					Once().Return(happyPathShareCurrency, nil)
			},
			expectedCode:  http.StatusOK,
			expectedBody:  happyPathShareCurrencyInBytes,
			assertNoCalls: false,
		},
		{
			name:          "Err: Auth err",
			token:         "",
			request:       service.ShareCurrencyByRequest{Figi: "BBG007N0Z367"},
			setupMock:     func(is *mocks.InstrumentService, token string, uid string) {},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"error": "Incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:          "Err: Invalid request",
			token:         tokens.OnlyReadToken,
			request:       `{"invalid request"}`,
			setupMock:     func(is *mocks.InstrumentService, token string, uid string) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "Invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:    "Err: GetClient err ",
			token:   tokens.OnlyReadToken,
			request: service.ShareCurrencyByRequest{Figi: "BBG007N0Z367"},
			setupMock: func(is *mocks.InstrumentService, token string, uid string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.InstrumentServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:    "Err: GetGetCurrencyBy err",
			token:   tokens.OnlyReadToken,
			request: service.ShareCurrencyByRequest{Figi: "BBG007N0Z367"},
			setupMock: func(is *mocks.InstrumentService, token string, figi string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetShareCurrencyBy", figi).
					Once().
					Return(service.ShareCurrencyByResponse{}, errors.New("op:service.GetShareCurrencyBy, error: can't get share by figi"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"error": "Could not get currency"}`,
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			analyticMock := mocks.NewAnalyticsService(t)
			portfolioMock := mocks.NewPortfolioService(t)
			instrumentMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMock(instrumentMock, tc.token, "")
			case []byte:
				jsonData = body
				tc.setupMock(instrumentMock, tc.token, "")
			case service.ShareCurrencyByRequest:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, body.Figi)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, "")
			}

			serviceMocks := service.NewService(analyticMock, portfolioMock, instrumentMock)
			handler := NewHandlers(serviceMocks)

			c, rec := createTestContext(tc.token, http.MethodPost, "/share/currecny", jsonData)

			err := handler.GetShareCurrencyBy(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				analyticMock.AssertNotCalled(t, "GetClient")
				analyticMock.AssertNotCalled(t, "GetShareCurrencyBy")
			} else {
				analyticMock.AssertExpectations(t)
			}

		})
	}
}

func TestGetFutureBy(t *testing.T) {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)

	cases := []struct {
		name          string
		token         string
		request       any
		setupMock     func(*mocks.InstrumentService, string, string)
		expectedCode  int
		expectedBody  string
		assertNoCalls bool
	}{
		{
			name:    "HappyPath",
			token:   tokens.OnlyReadToken,
			request: service.FutureReq{Figi: "FUTCNY032300"},
			setupMock: func(is *mocks.InstrumentService, token string, figi string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetFutureBy", figi).
					Once().Return(happyPathFutureBy, nil)
			},
			expectedCode:  http.StatusOK,
			expectedBody:  happyPathFutureByInBytes,
			assertNoCalls: false,
		},
		{
			name:          "Err: Auth err",
			token:         "",
			request:       service.FutureReq{Figi: "FUTCNY032300"},
			setupMock:     func(is *mocks.InstrumentService, token string, figi string) {},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"error": "Incorrect token"}`,
			assertNoCalls: true,
		},
		{
			name:          "Err: Invalid request",
			token:         tokens.OnlyReadToken,
			request:       `{"invalid request"}`,
			setupMock:     func(is *mocks.InstrumentService, token string, figi string) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "Invalid request"}`,
			assertNoCalls: true,
		},
		{
			name:    "Err: GetClient err ",
			token:   tokens.OnlyReadToken,
			request: service.FutureReq{Figi: "FUTCNY032300"},
			setupMock: func(is *mocks.InstrumentService, token string, figi string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(errors.New("op:sevrice.InstrumentServiceClient.GetClient, err: can't connect with tinkoffApi client"))
			},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"error": "TinkoffApi does not accesept token"}`,
			assertNoCalls: false,
		},
		{
			name:    "Err: GetGetCurrencyBy err",
			token:   tokens.OnlyReadToken,
			request: service.FutureReq{Figi: "FUTCNY032300"},
			setupMock: func(is *mocks.InstrumentService, token string, figi string) {
				is.On("GetClient", mock.Anything, token).
					Once().
					Return(nil)

				is.On("GetFutureBy", figi).
					Once().
					Return(service.Future{}, errors.New("op:service.GetFutureBy, error: can't get futures by figi"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"error": "Could not get futures"}`,
			assertNoCalls: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			analyticMock := mocks.NewAnalyticsService(t)
			portfolioMock := mocks.NewPortfolioService(t)
			instrumentMock := mocks.NewInstrumentService(t)
			var jsonData []byte
			switch body := tc.request.(type) {
			case string:
				jsonData = []byte(body)
				tc.setupMock(instrumentMock, tc.token, "")
			case []byte:
				jsonData = body
				tc.setupMock(instrumentMock, tc.token, "")
			case service.FutureReq:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, body.Figi)
			default:
				jsonData, _ = json.Marshal(body)
				tc.setupMock(instrumentMock, tc.token, "")
			}

			serviceMocks := service.NewService(analyticMock, portfolioMock, instrumentMock)
			handler := NewHandlers(serviceMocks)

			c, rec := createTestContext(tc.token, http.MethodPost, "/future", jsonData)

			err := handler.GetFutureBy(c)

			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
			if tc.assertNoCalls {
				analyticMock.AssertNotCalled(t, "GetClient")
				analyticMock.AssertNotCalled(t, "GetFutureBy")
			} else {
				analyticMock.AssertExpectations(t)
			}

		})
	}
}

func createTestContextWithoutBody(token, method, path string) (echo.Context, *httptest.ResponseRecorder) {
	var req *http.Request

	req = httptest.NewRequest(method, path, nil)

	e := echo.New()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func createTestContextFoAuth(token, method, path, body, prefix string, wantHeader bool) (echo.Context, *httptest.ResponseRecorder) {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	}
	e := echo.New()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if wantHeader {
		req.Header.Set("Authorization", prefix+" "+token)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func createTestContext(token, method, path string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))

	e := echo.New()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
