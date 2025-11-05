package service

import (
	"context"
	"testing"
	"time"
	"tinkoffApi/internal/configs"
	"tinkoffApi/lib/e/logger/loggerdicard"
	testhelpfunc "tinkoffApi/lib/testHelpFunc"
	"tinkoffApi/pkg/app"

	"github.com/stretchr/testify/require"
)

func TestGetClient_AnalyticService(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			wantErr: false,
		},
		{
			name:    "Unauthenticated desc = 40003",
			token:   "ijfoerigtj[o]",
			wantErr: true,
		},
		{
			name:    "Sandbox token",
			token:   tokens.SandboxToken,
			wantErr: false, //Как оказалось на практике
		},
		{
			name:    "Delete token err",
			token:   tokens.DeleteToken,
			wantErr: true,
		},
		{
			name:    "One account read token",
			token:   tokens.OneAccountReadToken,
			wantErr: false,
		},
		{
			name:    "All acsess token",
			token:   tokens.AllAcsessToken,
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			wantErr: false,
		},
		{
			name:    "Err: empty token",
			token:   "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewAnalyticsServiceClient(cnfgs.TinkoffApiConfig, logg)
			err := client.GetClient(ctx, tc.token)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})

	}
}

func TestGetClient_InstrumentService(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			wantErr: false,
		},
		{
			name:    "Unauthenticated desc = 40003",
			token:   "ijfoerigtj[o]",
			wantErr: true,
		},
		{
			name:    "Sandbox token",
			token:   tokens.SandboxToken,
			wantErr: false, //Как оказалось на практике
		},
		{
			name:    "Delete token err",
			token:   tokens.DeleteToken,
			wantErr: true,
		},
		{
			name:    "One account read token",
			token:   tokens.OneAccountReadToken,
			wantErr: false,
		},
		{
			name:    "All acsess token",
			token:   tokens.AllAcsessToken,
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			wantErr: false,
		},
		{
			name:    "Err: empty token",
			token:   "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewInstrumentServiceClient(cnfgs.TinkoffApiConfig, logg)
			err := client.GetClient(ctx, tc.token)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})

	}
}

func TestGetClient_PortfolioService(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			wantErr: false,
		},
		{
			name:    "Unauthenticated desc = 40003",
			token:   "ijfoerigtj[o]",
			wantErr: true,
		},
		{
			name:    "Sandbox token",
			token:   tokens.SandboxToken,
			wantErr: false, //Как оказалось на практике
		},
		{
			name:    "Delete token err",
			token:   tokens.DeleteToken,
			wantErr: true,
		},
		{
			name:    "One account read token",
			token:   tokens.OneAccountReadToken,
			wantErr: false,
		},
		{
			name:    "All acsess token",
			token:   tokens.AllAcsessToken,
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			wantErr: false,
		},
		{
			name:    "Err: empty token",
			token:   "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			err := client.GetClient(ctx, tc.token)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})

	}
}

func TestGetAccounts(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			wantErr: false,
		},
		{
			name:    "Sandbox token",
			token:   tokens.SandboxToken,
			wantErr: true, //Как оказалось на практике
		},
		{
			name:    "One account read token",
			token:   tokens.OneAccountReadToken,
			wantErr: false,
		},
		{
			name:    "All acsess token",
			token:   tokens.AllAcsessToken,
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			wantErr: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetAccounts()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetPortfolio(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		request PortfolioRequest
		wantErr bool
	}{
		{
			name:  "sucsees",
			token: tokens.OnlyReadToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 2,
			},
			wantErr: false,
		},
		{
			name:  "All acsess token",
			token: tokens.AllAcsessToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 2,
			},
			wantErr: false,
		},
		{
			name:  "Only trading token",
			token: tokens.OnlyTradingToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 2,
			},
			wantErr: false,
		},
		{
			name:  "Error: Close account",
			token: tokens.OnlyTradingToken,
			request: PortfolioRequest{
				AccountID:     "2012259491",
				AccountStatus: 3,
			},
			wantErr: true,
		},
		{
			name:    "Error: Sandbox token without request",
			token:   tokens.SandboxToken,
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Error: Sandbox token",
			token: tokens.SandboxToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 2,
			},
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Error: Token have not acsess to acount",
			token: tokens.OneAccountReadToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 2,
			},
			wantErr: true,
		},
		{
			name:  "Error: Empty accountID",
			token: tokens.OneAccountReadToken,
			request: PortfolioRequest{
				AccountID:     "",
				AccountStatus: 2,
			},
			wantErr: true,
		},
		{
			name:  "Error: Incorrect Status == 0 of accountId",
			token: tokens.OnlyReadToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 0,
			},
			wantErr: true,
		},
		{
			name:  "Error: Incorrect Status == 1 of accountId",
			token: tokens.OnlyReadToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 1,
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetPortfolio(tc.request)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

// func TestGetOperations_TimeNow(t *testing.T) {
// 	logg := loggerdicard.NewLoggerDiscard()
// 	app.MustInitialize()
// 	rootPath := app.MustGetRoot()
// 	tokens := testhelpfunc.MustTokensForTest(rootPath)
// 	cases := []struct {
// 		name    string
// 		token   string
// 		request OperationsRequest
// 		wantErr bool
// 	}{
// 		// Из-за рассинхрона времени с Тинькофф Апи. Тесты то работают, то падают.
// 		// Для обработки данной проблемы создал функцию MakeSafeGetOperationsRequest
// 		{
// 			name:  "sucsees",
// 			token: tokens.OnlyReadToken,
// 			request: OperationsRequest{
// 				AccountID: "2007907898",
// 				Date:      time.Now(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:  "All access token",
// 			token: tokens.AllAcsessToken,
// 			request: OperationsRequest{
// 				AccountID: "2007907898",
// 				Date:      time.Now(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:  "Only trading token",
// 			token: tokens.OnlyTradingToken,
// 			request: OperationsRequest{
// 				AccountID: "2007907898",
// 				Date:      time.Now(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:  "Close account",
// 			token: tokens.OnlyTradingToken,
// 			request: OperationsRequest{
// 				AccountID: "2012259491",
// 				Date:      time.Now(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "Sandbox token without request error",
// 			token:   tokens.SandboxToken,
// 			wantErr: true,
// 		},
// 		{
// 			name:  "Sandbox token error",
// 			token: tokens.SandboxToken,
// 			request: OperationsRequest{
// 				AccountID: "2007907898",
// 				Date:      time.Now(),
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:  "Err: Token have not acsess to acount",
// 			token: tokens.OneAccountReadToken,
// 			request: OperationsRequest{
// 				AccountID: "2007907898",
// 				Date:      time.Now(),
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:  "Err: empty accountId",
// 			token: tokens.OnlyReadToken,
// 			request: OperationsRequest{
// 				AccountID: "",
// 				Date:      time.Now(),
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			cnfgs := configs.MustInitConfigs(rootPath)
// 			cnfgs.TinkoffApiConfig.Token = tc.token
// 			ctx := context.Background()

//          client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
//          client.GetClient(ctx, tc.token)
// 			_, err := client.GetOperations(tc.request)
// 			if tc.wantErr {
// 				require.Error(t, err)
// 				return
// 			}
// 			require.NoError(t, err)

// 		})
// 	}
// }

func TestGetOperations_TimeFromFuture(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		request OperationsRequest
		wantErr bool
	}{
		{
			name:  "error",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "All access token error",
			token: tokens.AllAcsessToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Only trading token error",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Close account error",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2012259491",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:    "Sandbox token without request error",
			token:   tokens.SandboxToken,
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Sandbox token error",
			token: tokens.SandboxToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Err: Token have not acsess to acount",
			token: tokens.OneAccountReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Err: empty accountId",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "",
				Date:      time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetOperations(tc.request)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetOperations_TimeFromPast(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		request OperationsRequest
		wantErr bool
	}{
		{
			name:  "sucsees",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:  "All access token",
			token: tokens.AllAcsessToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:  "Only trading token",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:  "Close account",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2012259491",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:    "Sandbox token without request error",
			token:   tokens.SandboxToken,
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Sandbox token error",
			token: tokens.SandboxToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Err: Token have not acsess to acount",
			token: tokens.OneAccountReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Err: empty accountId",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "",
				Date:      time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetOperations(tc.request)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestMakeSafeGetOperationsRequest_TimeNow(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		request OperationsRequest
		wantErr bool
	}{
		{
			name:  "sucsees",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now(),
			},
			wantErr: false,
		},
		{
			name:  "All access token",
			token: tokens.AllAcsessToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now(),
			},
			wantErr: false,
		},
		{
			name:  "Only trading token",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now(),
			},
			wantErr: false,
		},
		{
			name:  "Close account",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2012259491",
				Date:      time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "Sandbox token without request error",
			token:   tokens.SandboxToken,
			wantErr: true,
		},
		{
			name:  "Sandbox token error",
			token: tokens.SandboxToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now(),
			},
			wantErr: true,
		},
		{
			name:  "Err: Token have not acsess to acount",
			token: tokens.OneAccountReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now(),
			},
			wantErr: true,
		},
		{
			name:  "Err: empty accountId",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "",
				Date:      time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.MakeSafeGetOperationsRequest(tc.request)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
func TestMakeSafeGetOperationsRequest_TimeFromFuture(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		request OperationsRequest
		wantErr bool
	}{
		{
			name:  "error",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "All access token error",
			token: tokens.AllAcsessToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Only trading token error",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Close account error",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2012259491",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:    "Sandbox token without request error",
			token:   tokens.SandboxToken,
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Sandbox token error",
			token: tokens.SandboxToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Err: Token have not acsess to acount",
			token: tokens.OneAccountReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(10, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Err: empty accountId",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "",
				Date:      time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.MakeSafeGetOperationsRequest(tc.request)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestMakeSafeGetOperationsRequest_TimeFromPast(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		request OperationsRequest
		wantErr bool
	}{
		{
			name:  "sucsees",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:  "All access token",
			token: tokens.AllAcsessToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:  "Only trading token",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:  "Close account",
			token: tokens.OnlyTradingToken,
			request: OperationsRequest{
				AccountID: "2012259491",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: false,
		},
		{
			name:    "Sandbox token without request error",
			token:   tokens.SandboxToken,
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Sandbox token error",
			token: tokens.SandboxToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Err: Token have not acsess to acount",
			token: tokens.OneAccountReadToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now().AddDate(-100, 0, 0),
			},
			wantErr: true,
		},
		{
			name:  "Err: empty accountId",
			token: tokens.OnlyReadToken,
			request: OperationsRequest{
				AccountID: "",
				Date:      time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.MakeSafeGetOperationsRequest(tc.request)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
func TestAllAssetUids(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			wantErr: false,
		},
		{
			name:    "All access token",
			token:   tokens.AllAcsessToken,
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			wantErr: false,
		},
		{
			name:    "Close account",
			token:   tokens.OnlyTradingToken,
			wantErr: false,
		},
		{
			name:    "Sandbox token without request error",
			token:   tokens.SandboxToken,
			wantErr: false, //Как оказалось на практике
		},
		{
			name:    "Sandbox token error",
			token:   tokens.SandboxToken,
			wantErr: false, //Как оказалось на практике
		},
		{
			name:    "Err: Token have not acsess to acount",
			token:   tokens.OneAccountReadToken,
			wantErr: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewAnalyticsServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetAllAssetUids()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetFutureBy(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		figi    string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			figi:    "FUTCNY032300",
			wantErr: false,
		},
		{
			name:    "All access token",
			token:   tokens.AllAcsessToken,
			figi:    "FUTCNY032300",
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			figi:    "FUTCNY032300",
			wantErr: false,
		},
		{
			name:    "Close account",
			token:   tokens.OnlyTradingToken,
			figi:    "FUTCNY032300",
			wantErr: false,
		},
		{
			name:  "Sandbox token without request error",
			token: tokens.SandboxToken,

			wantErr: true,
		},
		{
			name:    "Sandbox token ",
			token:   tokens.SandboxToken,
			figi:    "FUTCNY032300",
			wantErr: false,
		},
		{
			name:    "Token have not acsess to acount",
			token:   tokens.OneAccountReadToken,
			figi:    "FUTCNY032300",
			wantErr: false,
		},
		{
			name:    "Err: Incorrect figi",
			token:   tokens.OneAccountReadToken,
			figi:    "FUTCNY03ghgerhgrehrt2300",
			wantErr: true,
		},
		{
			name:    "Err: empty figi",
			token:   tokens.OnlyReadToken,
			figi:    "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewInstrumentServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetFutureBy(tc.figi)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetBondByUid(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		uid     string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "All access token",
			token:   tokens.AllAcsessToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Close account",
			token:   tokens.OnlyTradingToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:  "Sandbox token without request error",
			token: tokens.SandboxToken,

			wantErr: true,
		},
		{
			name:    "Sandbox token ",
			token:   tokens.SandboxToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Token have not acsess to acount",
			token:   tokens.OneAccountReadToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Err: Incorrect uid",
			token:   tokens.OneAccountReadToken,
			uid:     "FUTCNY03ghgerhgrehrt2300",
			wantErr: true,
		},
		{
			name:    "Empty uid",
			token:   tokens.OnlyReadToken,
			uid:     "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewInstrumentServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetBondByUid(tc.uid)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetCurrencyBy(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		figi    string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			figi:    "CNY000TODTOM",
			wantErr: false,
		},
		{
			name:    "All access token",
			token:   tokens.AllAcsessToken,
			figi:    "CNY000TODTOM",
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			figi:    "CNY000TODTOM",
			wantErr: false,
		},
		{
			name:    "Close account",
			token:   tokens.OnlyTradingToken,
			figi:    "CNY000TODTOM",
			wantErr: false,
		},
		{
			name:  "Sandbox token without request error",
			token: tokens.SandboxToken,

			wantErr: true,
		},
		{
			name:    "Sandbox token",
			token:   tokens.SandboxToken,
			figi:    "CNY000TODTOM",
			wantErr: false,
		},
		{
			name:    "Token have not acsess to acount",
			token:   tokens.OneAccountReadToken,
			figi:    "CNY000TODTOM",
			wantErr: false,
		},
		{
			name:    "Err: Incorrect figi",
			token:   tokens.OneAccountReadToken,
			figi:    "FUTCNY03ghgerhgrehrt2300",
			wantErr: true,
		},
		{
			name:    "Err: Futures figi",
			token:   tokens.OnlyReadToken,
			figi:    "FUTCNY032300",
			wantErr: true,
		},
		{
			name:    "Err: empty string",
			token:   tokens.OnlyReadToken,
			figi:    "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewInstrumentServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetCurrencyBy(tc.figi)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestFindBy(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		query   string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			query:   "e80d1280-d512-4755-b48b-1187fd6cb2d8",
			wantErr: false,
		},
		{
			name:    "All access token",
			token:   tokens.AllAcsessToken,
			query:   "e80d1280-d512-4755-b48b-1187fd6cb2d8",
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			query:   "e80d1280-d512-4755-b48b-1187fd6cb2d8",
			wantErr: false,
		},
		{
			name:    "Close account",
			token:   tokens.OnlyTradingToken,
			query:   "e80d1280-d512-4755-b48b-1187fd6cb2d8",
			wantErr: false,
		},
		{
			name:  "Sandbox token without request error",
			token: tokens.SandboxToken,

			wantErr: true,
		},
		{
			name:    "Sandbox token",
			token:   tokens.SandboxToken,
			query:   "e80d1280-d512-4755-b48b-1187fd6cb2d8",
			wantErr: false,
		},
		{
			name:    "Token have not acsess to acount",
			token:   tokens.OneAccountReadToken,
			query:   "e80d1280-d512-4755-b48b-1187fd6cb2d8",
			wantErr: false,
		},
		{
			name:    "Err: Incorrect query",
			token:   tokens.OneAccountReadToken,
			query:   "e80d1280-d512-4755-b48b-1187regerfd6cb2d8",
			wantErr: false,
		},
		{
			name:    "Err: Futures figi",
			token:   tokens.OnlyReadToken,
			query:   "FUTCNY032300",
			wantErr: false,
		},
		{
			name:    "Err: empty query",
			token:   tokens.OnlyReadToken,
			query:   "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewInstrumentServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.FindBy(tc.query)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetBondsActions(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		uid     string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "All access token",
			token:   tokens.AllAcsessToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Close account",
			token:   tokens.OnlyTradingToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:  "Sandbox token without request error",
			token: tokens.SandboxToken,

			wantErr: true,
		},
		{
			name:    "Sandbox token ",
			token:   tokens.SandboxToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Token have not acsess to acount",
			token:   tokens.OneAccountReadToken,
			uid:     "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr: false,
		},
		{
			name:    "Err: Incorrect uid",
			token:   tokens.OneAccountReadToken,
			uid:     "FUTCNY03ghgerhgrehrt2300",
			wantErr: true,
		},
		{
			name:    "Empty uid",
			token:   tokens.OnlyReadToken,
			uid:     "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewAnalyticsServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetBondsActions(tc.uid)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetLastPriceInPersentageToNominal(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name          string
		token         string
		instrumentUid string
		wantErr       bool
	}{
		{
			name:          "sucsees",
			token:         tokens.OnlyReadToken,
			instrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr:       false,
		},
		{
			name:          "All access token",
			token:         tokens.AllAcsessToken,
			instrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr:       false,
		},
		{
			name:          "Only trading token",
			token:         tokens.OnlyTradingToken,
			instrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr:       false,
		},
		{
			name:          "Close account",
			token:         tokens.OnlyTradingToken,
			instrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr:       false,
		},
		{
			name:  "Sandbox token without request error",
			token: tokens.SandboxToken,

			wantErr: true,
		},
		{
			name:          "Sandbox token ",
			token:         tokens.SandboxToken,
			instrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr:       false,
		},
		{
			name:          "Token have not acsess to acount",
			token:         tokens.OneAccountReadToken,
			instrumentUid: "070d82ad-e9e0-41e4-8eca-cbe9f5830db2",
			wantErr:       false,
		},
		{
			name:          "Err: Incorrect uid",
			token:         tokens.OneAccountReadToken,
			instrumentUid: "invalid_uid",
			wantErr:       true,
		},
		{
			name:          "Empty uid",
			token:         tokens.OnlyReadToken,
			instrumentUid: "",
			wantErr:       true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewAnalyticsServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetLastPriceInPersentageToNominal(tc.instrumentUid)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetShareCurrencyBy(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := testhelpfunc.MustTokensForTest(rootPath)
	cases := []struct {
		name    string
		token   string
		figi    string
		wantErr bool
	}{
		{
			name:    "sucsees",
			token:   tokens.OnlyReadToken,
			figi:    "BBG004S68FR6",
			wantErr: false,
		},
		{
			name:    "All access token",
			token:   tokens.AllAcsessToken,
			figi:    "BBG004S68FR6",
			wantErr: false,
		},
		{
			name:    "Only trading token",
			token:   tokens.OnlyTradingToken,
			figi:    "BBG004S68FR6",
			wantErr: false,
		},
		{
			name:    "Close account",
			token:   tokens.OnlyTradingToken,
			figi:    "BBG004S68FR6",
			wantErr: false,
		},
		{
			name:  "Sandbox token without request error",
			token: tokens.SandboxToken,

			wantErr: true,
		},
		{
			name:    "Sandbox token",
			token:   tokens.SandboxToken,
			figi:    "BBG004S68FR6",
			wantErr: false,
		},
		{
			name:    "Token have not acsess to acount",
			token:   tokens.OneAccountReadToken,
			figi:    "BBG004S68FR6",
			wantErr: false,
		},
		{
			name:    "Err: Incorrect figi",
			token:   tokens.OneAccountReadToken,
			figi:    "FUTCNY03ghgerhgrehrt2300",
			wantErr: true,
		},
		{
			name:    "Err: Futures figi",
			token:   tokens.OnlyReadToken,
			figi:    "FUTCNY032300",
			wantErr: true,
		},
		{
			name:    "Err: empty string",
			token:   tokens.OnlyReadToken,
			figi:    "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx := context.Background()

			client := NewInstrumentServiceClient(cnfgs.TinkoffApiConfig, logg)
			client.GetClient(ctx, tc.token)
			_, err := client.GetShareCurrencyBy(tc.figi)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
