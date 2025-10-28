package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"
	"tinkoffApi/internal/configs"
	"tinkoffApi/lib/e/logger/loggerdicard"
	"tinkoffApi/pkg/app"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

type Tokens struct {
	OnlyReadToken       string
	SandboxToken        string
	DeleteToken         string
	OneAccountReadToken string
	AllAcsessToken      string
	OnlyTradingToken    string
}

func MustTokensForTest(rootPath string) *Tokens {
	const op = "service.MustTokenForTest"
	envPath := filepath.Join(rootPath, "tokens.env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("%s:Could not find any .env files:%s", op, err)
	}

	var tokens Tokens
	token := os.Getenv("TEST_TINKOFF_TOKEN")
	if token == "" {
		log.Fatalf("%s: TEST_TINKOFF_TOKEN is not set", op)
	}
	tokens.OnlyReadToken = token

	token = os.Getenv("TEST_SANDBOX_TOKEN")
	if token == "" {
		log.Fatalf("%s: TEST_SANDBOX_TOKEN is not set", op)
	}
	tokens.SandboxToken = token

	token = os.Getenv("TEST_DELETE_TOKEN")
	if token == "" {
		log.Fatalf("%s: TEST_DELETE_TOKEN is not set", op)
	}

	tokens.DeleteToken = token

	token = os.Getenv("TEST_ONE_ACOUNT_TOKEN")
	if token == "" {
		log.Fatalf("%s: TEST_ONE_ACOUNT_TOKEN is not set", op)
	}

	tokens.OneAccountReadToken = token

	token = os.Getenv("TEST_ALL_ACSESS_TOKEN")
	if token == "" {
		log.Fatalf("%s: TEST_ALL_ACSESS_TOKEN is not set", op)
	}

	tokens.AllAcsessToken = token

	token = os.Getenv("TEST_ONLY_TRADING_TOKEB")
	if token == "" {
		log.Fatalf("%s: TEST_ONLY_TRADING_TOKEB is not set", op)
	}

	tokens.OnlyTradingToken = token

	return &tokens

}

func TestGetClient(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := MustTokensForTest(rootPath)
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token

			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer cancel()
			client := New(ctx, logg, cnfgs.TinkoffApiConfig)
			err := client.getClient()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestFillCLient(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := MustTokensForTest(rootPath)
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer cancel()
			client := New(ctx, logg, cnfgs.TinkoffApiConfig)
			err := client.FillClient(tc.token)
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
	tokens := MustTokensForTest(rootPath)
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
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer cancel()
			client := New(ctx, logg, cnfgs.TinkoffApiConfig)
			client.FillClient(tc.token)
			_, err := client.GetAccounts()
			fmt.Println(err)
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
	tokens := MustTokensForTest(rootPath)
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
			name:  "Close account error",
			token: tokens.OnlyTradingToken,
			request: PortfolioRequest{
				AccountID:     "2012259491",
				AccountStatus: 3,
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
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 2,
			},
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Err: Token have not acsess to acount",
			token: tokens.OneAccountReadToken,
			request: PortfolioRequest{
				AccountID:     "2007907898",
				AccountStatus: 2,
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer cancel()
			client := New(ctx, logg, cnfgs.TinkoffApiConfig)
			client.FillClient(tc.token)
			_, err := client.GetPortfolio(tc.request)
			fmt.Println(err)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetOperations_TimeNow(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := MustTokensForTest(rootPath)
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
			wantErr: true, //Как оказалось на практике
		},
		{
			name:  "Sandbox token error",
			token: tokens.SandboxToken,
			request: OperationsRequest{
				AccountID: "2007907898",
				Date:      time.Now(),
			},
			wantErr: true, //Как оказалось на практике
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer cancel()
			client := New(ctx, logg, cnfgs.TinkoffApiConfig)
			client.FillClient(tc.token)
			_, err := client.GetOperations(tc.request)
			fmt.Println(err)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestGetOperations_TimeFromFuture(t *testing.T) {
	logg := loggerdicard.NewLoggerDiscard()
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	tokens := MustTokensForTest(rootPath)
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer cancel()
			client := New(ctx, logg, cnfgs.TinkoffApiConfig)
			client.FillClient(tc.token)
			_, err := client.GetOperations(tc.request)
			fmt.Println(err)
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
	tokens := MustTokensForTest(rootPath)
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cnfgs := configs.MustInitConfigs(rootPath)
			cnfgs.TinkoffApiConfig.Token = tc.token
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer cancel()
			client := New(ctx, logg, cnfgs.TinkoffApiConfig)
			client.FillClient(tc.token)
			_, err := client.GetOperations(tc.request)
			fmt.Println(err)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
