package testhelpfunc

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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
