package hanlders

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"tinkoffApi/internal/service"
	"tinkoffApi/lib/cryptoToken"
	"tinkoffApi/lib/valuefromcontext"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	service      *service.Service
	tokenCrypter *cryptoToken.TokenCrypter
	redis        *redis.Client
}

func NewHandlers(service *service.Service, tokenCrypter *cryptoToken.TokenCrypter, redis *redis.Client) *Handlers {
	return &Handlers{
		service:      service,
		tokenCrypter: tokenCrypter,
		redis:        redis,
	}
}

var errHeaderRequierd error = errors.New("header auth requierd")
var errInvalidAuthFormat error = errors.New("invalid Authorization format, expected: Bearer <token>")
var errEmptyToken error = errors.New("empty token")
var errNoTokenInRedis error = errors.New("token not found for the provided user id")
var errRedisDoNotAnswer error = errors.New("token storage temporarily unavailable")

func (h *Handlers) AuthCheckTokenMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		chatID := c.Request().Header.Get(HeaderChatID)

		if chatID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": errHeaderRequierd.Error()})
		}
		ctx := c.Request().Context()
		tokenInBase64, err := h.redis.Get(ctx, chatID).Result()
		if err == redis.Nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": errNoTokenInRedis.Error()})
		}
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": errRedisDoNotAnswer.Error()})
		}
		decodedJson, err := base64.StdEncoding.DecodeString(tokenInBase64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errHeaderRequierd.Error()})
		}
		var encrypredToken cryptoToken.EncryptedToken
		err = json.Unmarshal(decodedJson, &encrypredToken)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errHeaderRequierd.Error()})
		}
		token, err := cryptoToken.DecryptToken(&encrypredToken, h.tokenCrypter.Key)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errInvalidAuthFormat.Error()})
		}
		ctx = context.WithValue(c.Request().Context(), valuefromcontext.EncryptedTokenKey, token)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func (h *Handlers) CheckToken(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := h.service.PortfolioService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	_, err = h.service.PortfolioService.GetAccounts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get accounts"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handlers) GetAccounts(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := h.service.PortfolioService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	accs, err := h.service.PortfolioService.GetAccounts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get accounts"})
	}
	return c.JSON(http.StatusOK, accs)
}

func (h *Handlers) GetPortfolio(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var portffolioReq service.PortfolioRequest
	err := c.Bind(&portffolioReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err = h.service.PortfolioService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}

	portf, err := h.service.PortfolioService.GetPortfolio(portffolioReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get portfolio"})
	}
	return c.JSON(http.StatusOK, portf)
}

func (h *Handlers) GetOperations(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var operationReq service.OperationsRequest
	err := c.Bind(&operationReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err = h.service.PortfolioService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}

	operations, err := h.service.PortfolioService.MakeSafeGetOperationsRequest(operationReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get operations"})
	}
	return c.JSON(http.StatusOK, operations)
}

func (h *Handlers) GetAllAssetUids(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := h.service.AnalyticsService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	allAssetUids, err := h.service.AnalyticsService.GetAllAssetUids()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get assets uids"})
	}
	return c.JSON(http.StatusOK, allAssetUids)
}

func (h *Handlers) GetFutureBy(c echo.Context) error {

	var body service.FutureReq
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = h.service.InstrumentService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	future, err := h.service.InstrumentService.GetFutureBy(body.Figi)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get futures"})
	}
	return c.JSON(http.StatusOK, future)
}

func (h *Handlers) GetBondBy(c echo.Context) error {
	var body service.BondReq
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = h.service.InstrumentService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	bond, err := h.service.InstrumentService.GetBondByUid(body.Uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get bond"})
	}
	return c.JSON(http.StatusOK, bond)
}

func (h *Handlers) GetCurrencyBy(c echo.Context) error {
	var body service.CurrencyReq
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = h.service.InstrumentService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	currency, err := h.service.InstrumentService.GetCurrencyBy(body.Figi)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get currecny"})
	}
	return c.JSON(http.StatusOK, currency)
}

func (h *Handlers) GetShareCurrencyBy(c echo.Context) error {
	var body service.ShareCurrencyByRequest
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = h.service.InstrumentService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	currency, err := h.service.InstrumentService.GetShareCurrencyBy(body.Figi)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get currency"})
	}
	return c.JSON(http.StatusOK, currency)
}

func (h *Handlers) FindBy(c echo.Context) error {
	var body service.FindByReq
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = h.service.InstrumentService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	instruments, err := h.service.InstrumentService.FindBy(body.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get instruments"})
	}
	return c.JSON(http.StatusOK, instruments)
}

func (h *Handlers) GetBondsActions(c echo.Context) error {
	var body service.BondsActionsReq
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = h.service.AnalyticsService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	bondIdentificators, err := h.service.AnalyticsService.GetBondsActions(body.InstrumentUid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get bond identificators"})
	}
	return c.JSON(http.StatusOK, bondIdentificators)
}

func (h *Handlers) GetLastPriceInPersentageToNominal(c echo.Context) error {

	var body service.LastPriceReq
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = h.service.AnalyticsService.GetClient(ctx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "TinkoffApi does not accesept token"})
	}
	lastPrice, err := h.service.AnalyticsService.GetLastPriceInPersentageToNominal(body.InstrumentUid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get last price"})
	}

	return c.JSON(http.StatusOK, lastPrice)
}
