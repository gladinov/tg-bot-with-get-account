package hanlders

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tinkoffApi/internal/service"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	service service.Service
}

func NewHandlers(service service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

var errHeaderRequierd error = errors.New("header auth requierd")
var errInvalidAuthFormat error = errors.New("invalid Authorization format, expected: Bearer <token>")
var errEmptyToken error = errors.New("empty token")

func auth(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errHeaderRequierd
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errInvalidAuthFormat
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errEmptyToken
	}
	return token, nil
}

func (h *Handlers) GetAccounts(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	accs, err := h.service.GetAccounts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get accounts"})
	}
	return c.JSON(http.StatusOK, accs)
}

func (h *Handlers) GetPortfolio(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	var portffolioReq service.PortfolioRequest
	err = c.Bind(&portffolioReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}

	portf, err := h.service.GetPortfolio(portffolioReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get portfolio"})
	}
	return c.JSON(http.StatusOK, portf)

}

func (h *Handlers) GetOperations(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	var operationReq service.OperationsRequest
	err = c.Bind(&operationReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}

	operations, err := h.service.MakeSafeGetOperationsRequest(operationReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get operations"})
	}
	return c.JSON(http.StatusOK, operations)
}

func (h *Handlers) GetAllAssetUids(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	allAssetUids, err := h.service.GetAllAssetUids()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get assets uids"})
	}
	return c.JSON(http.StatusOK, allAssetUids)
}

func (h *Handlers) GetFutureBy(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	var body service.FutureReq
	err = c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	future, err := h.service.GetFutureBy(body.Figi)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not get futures, delete this block. err :%s", err.Error())})
	}
	return c.JSON(http.StatusOK, future)
}

func (h *Handlers) GetBondBy(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	var body service.BondReq
	err = c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	bond, err := h.service.GetBondByUid(body.Uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get bond"})
	}
	return c.JSON(http.StatusOK, bond)
}

func (h *Handlers) GetCurrencyBy(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	var body service.CurrencyReq
	err = c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	currency, err := h.service.GetCurrencyBy(body.Figi)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get currecny"})
	}
	return c.JSON(http.StatusOK, currency)
}

func (h *Handlers) GetBaseShareFutureValute(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	var body service.BaseShareFutureValuteReq
	err = c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	currency, err := h.service.GetBaseShareFutureValute(body.SharePositionUid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get currency "})
	}
	return c.JSON(http.StatusOK, currency)
}

func (h *Handlers) FindBy(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	var body service.FindByReq
	err = c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	instruments, err := h.service.FindBy(body.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get instruments"})
	}
	return c.JSON(http.StatusOK, instruments)
}

func (h *Handlers) GetBondsActions(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	var body service.BondsActionsReq
	err = c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	bondIdentificators, err := h.service.GetBondsActions(body.InstrumentUid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get bond identificators"})
	}
	return c.JSON(http.StatusOK, bondIdentificators)
}

func (h *Handlers) GetLastPriceInPersentageToNominal(c echo.Context) error {
	authHeader, err := auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	var body service.LastPriceReq
	err = c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	err = h.service.FillClient(authHeader)
	// TODO: Проверить тестами корректность токена! Если тут проходит проверка, то перенести FillClient в Auth
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "incorrect token"})
	}
	lastPrice, err := h.service.GetLastPriceInPersentageToNominal(body.InstrumentUid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get last price"})
	}

	return c.JSON(http.StatusOK, lastPrice)
}
