package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gladinov/valuefromcontext"

	"github.com/gladinov/e"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type InstrumentsServiceClient struct {
	Logg   investgo.Logger
	config *investgo.Config
}

type PortfolioServiceClient struct {
	Logg   investgo.Logger
	config *investgo.Config
}
type AnalyticsServiceClient struct {
	Logg   investgo.Logger
	config *investgo.Config
}

func NewInstrumentServiceClient(config *investgo.Config, logg investgo.Logger) InstrumentService {
	return &InstrumentsServiceClient{
		config: config,
		Logg:   logg,
	}
}

func NewPortfolioServiceClient(config *investgo.Config, logg investgo.Logger) PortfolioService {
	return &PortfolioServiceClient{
		config: config,
		Logg:   logg,
	}
}

func NewAnalyticsServiceClient(config *investgo.Config, logg investgo.Logger) AnalyticsService {
	return &AnalyticsServiceClient{
		config: config,
		Logg:   logg,
	}
}

func NewService(
	analyticsService AnalyticsService,
	portfolioService PortfolioService,
	instrumentService InstrumentService,
) *Service {
	return &Service{
		AnalyticsService:  analyticsService,
		PortfolioService:  portfolioService,
		InstrumentService: instrumentService,
	}
}

func (c *InstrumentsServiceClient) GetClient(ctx context.Context) (_ *investgo.Client, err error) {
	const op = "sevrice.InstrumentsServiceClient.GetClient"
	token, err := valuefromcontext.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	newConfig := *c.config
	newConfig.Token = token

	client, err := investgo.NewClient(ctx, newConfig, c.Logg)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err: can't connect with tinkoffApi client", op)
	}

	return client, nil
}

func (c *AnalyticsServiceClient) GetClient(ctx context.Context) (_ *investgo.Client, err error) {
	const op = "sevrice.AnalyticsServiceClient.GetClient"

	token, err := valuefromcontext.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	newConfig := *c.config
	newConfig.Token = token

	client, err := investgo.NewClient(ctx, newConfig, c.Logg)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err: can't connect with tinkoffApi client", op)
	}

	return client, nil
}

func (c *PortfolioServiceClient) GetClient(ctx context.Context) (_ *investgo.Client, err error) {
	const op = "sevrice.PortfolioServiceClient.GetClient"

	token, err := valuefromcontext.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	newConfig := *c.config
	newConfig.Token = token

	client, err := investgo.NewClient(ctx, newConfig, c.Logg)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err: can't connect with tinkoffApi client", op)
	}

	return client, nil
}

func (c *PortfolioServiceClient) GetAccounts(client *investgo.Client) (map[string]Account, error) {
	const op = "sevrice.GetAccounts"
	usersService := client.NewUsersServiceClient()
	accounts := make(map[string]Account)
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err: could not get accounts", op)
	}
	accs := accsResp.GetAccounts()
	for _, acc := range accs {
		account := Account{
			Id:          acc.GetId(),
			Type:        acc.GetType().String(),
			Name:        acc.GetName(),
			OpenedDate:  acc.GetOpenedDate().AsTime(),
			ClosedDate:  acc.GetClosedDate().AsTime(),
			Status:      int64(acc.GetStatus()),
			AccessLevel: int64(acc.GetAccessLevel()),
		}
		accounts[acc.GetId()] = account
	}

	return accounts, nil
}

func (c *PortfolioServiceClient) GetPortfolio(client *investgo.Client, request PortfolioRequest) (_ Portfolio, err error) {
	const op = "service.GetPortfolio"
	accountID := request.AccountID
	accountStatus := request.AccountStatus
	portfolio := Portfolio{}
	switch accountStatus {
	case 0:
		return Portfolio{}, ErrCloseAccount
	case 1:
		return Portfolio{}, ErrUnspecifiedAccount
	case 3:
		return Portfolio{}, ErrNewNotOpenYetAccount
	}
	if accountID == "" {
		return Portfolio{}, ErrEmptyAccountIdInRequest
	}
	operationsService := client.NewOperationsServiceClient()

	portfolioResp, err := operationsService.GetPortfolio(accountID,
		pb.PortfolioRequest_RUB)
	if err != nil {
		return Portfolio{}, fmt.Errorf("op: %s, error: can't get portifolio positions from tinkoff Api", op)
	}
	portfolio.Positions = ConvertPbToPortfolioPositions(portfolioResp.GetPositions())
	portfolio.TotalAmount = ConvertPbToMoneyValue(portfolioResp.GetTotalAmountPortfolio())

	return portfolio, nil
}

func ConvertPbToPortfolioPositions(pbPositions []*pb.PortfolioPosition) []PortfolioPositions {
	positions := make([]PortfolioPositions, 0, len(pbPositions))
	for _, pbPosition := range pbPositions {
		position := PortfolioPositions{
			Figi:                     pbPosition.GetFigi(),
			InstrumentType:           pbPosition.GetInstrumentType(),
			Quantity:                 ConvertPbToQuatation(pbPosition.GetQuantity()),
			AveragePositionPrice:     ConvertPbToMoneyValue(pbPosition.GetAveragePositionPrice()),
			ExpectedYield:            ConvertPbToQuatation(pbPosition.GetExpectedYield()),
			CurrentNkd:               ConvertPbToMoneyValue(pbPosition.GetCurrentNkd()),
			CurrentPrice:             ConvertPbToMoneyValue(pbPosition.GetCurrentPrice()),
			AveragePositionPriceFifo: ConvertPbToMoneyValue(pbPosition.GetAveragePositionPriceFifo()),
			Blocked:                  pbPosition.GetBlocked(),
			BlockedLots:              ConvertPbToQuatation(pbPosition.GetBlockedLots()),
			PositionUid:              pbPosition.GetPositionUid(),
			InstrumentUid:            pbPosition.GetInstrumentUid(),
			VarMargin:                ConvertPbToMoneyValue(pbPosition.GetVarMargin()),
			ExpectedYieldFifo:        ConvertPbToQuatation(pbPosition.GetExpectedYieldFifo()),
			DailyYield:               ConvertPbToMoneyValue(pbPosition.GetDailyYield()),
			Ticker:                   pbPosition.GetTicker(),
		}
		positions = append(positions, position)
	}
	return positions
}

func ConvertPbToMoneyValue(pbMoneyValue *pb.MoneyValue) MoneyValue {
	if pbMoneyValue == nil {
		return MoneyValue{}
	}
	return MoneyValue{
		Currency: pbMoneyValue.GetCurrency(),
		Units:    pbMoneyValue.GetUnits(),
		Nano:     pbMoneyValue.GetNano(),
	}
}

func ConvertPbToQuatation(pbQuatation *pb.Quotation) Quotation {
	if pbQuatation == nil {
		return Quotation{}
	}
	return Quotation{
		Units: pbQuatation.GetUnits(),
		Nano:  pbQuatation.GetNano(),
	}
}

func (c *PortfolioServiceClient) GetOperations(client *investgo.Client, request OperationsRequest) (_ []Operation, err error) {
	const op = "tinkoffApi.GetOperations"
	defer func() { err = e.WrapIfErr(op+": can't get operations from tinkoff API", err) }()

	accountID := request.AccountID
	if request.AccountID == "" {
		return nil, fmt.Errorf("op: %s, error: empty account ID", op)
	}
	date := request.Date.UTC()
	now := time.Now().UTC()

	if date.After(now) {
		return nil, fmt.Errorf("op:%s, from can't be more than the current date", op)
	}
	allOperations := make([]*pb.OperationItem, 0)

	operationsService := client.NewOperationsServiceClient()
	operationsResp, err := operationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
		AccountId: accountID,
		From:      date,
		To:        now,
		Limit:     1000,
	})
	if err != nil {
		return nil, fmt.Errorf("op:%s, failed to get operations: %w", op, err)
	}
	operations := operationsResp.GetOperationsByCursorResponse.GetItems()
	allOperations = append(allOperations, operations...)
	nextCursor := operationsResp.NextCursor
	for nextCursor != "" {
		operationsResp, err := operationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
			AccountId: accountID,
			Limit:     1000,
			Cursor:    nextCursor,
		})
		if err != nil {
			return nil, fmt.Errorf("op:%s, failed to get operations with cursor: %w", op, err)
		}
		nextCursor = operationsResp.NextCursor
		operations := operationsResp.GetOperationsByCursorResponse.Items
		allOperations = append(allOperations, operations...)

		// TODO: Refactor this block. Add more context.
		if len(allOperations) > 10000 {
			break
		}

	}
	resp := convertOperationsPbToOperaions(allOperations)
	return resp, nil
}

func (c *PortfolioServiceClient) MakeSafeGetOperationsRequest(client *investgo.Client, request OperationsRequest) ([]Operation, error) {
	var lastErr error

	// Пробуем с разными сдвигами времени
	for _, offset := range []time.Duration{
		0,
		-1 * time.Minute,
		-2 * time.Minute,
		-3 * time.Minute,
	} {
		adjustedRequest := adjustRequestTime(request, offset)
		operations, err := c.GetOperations(client, adjustedRequest)
		if err == nil {
			return operations, nil
		}

		if !isTimeError(err) {
			return nil, err
		}

		lastErr = err
	}

	return nil, lastErr
}

func isTimeError(inputErr error) bool {
	if inputErr == nil {
		return false
	}

	errorStr := inputErr.Error()
	return strings.Contains(errorStr, "30070")
}

func adjustRequestTime(request OperationsRequest, offset time.Duration) OperationsRequest {
	adjustedRequest := request
	if !request.Date.IsZero() {
		adjustedRequest.Date = request.Date.Add(offset).UTC()
	}
	return adjustedRequest
}

func convertOperationsPbToOperaions(operations []*pb.OperationItem) []Operation {
	transformOperations := make([]Operation, 0, len(operations))
	for _, v := range operations {
		transformOperation := Operation{
			Currency:          v.GetPrice().Currency,
			BrokerAccountId:   v.GetBrokerAccountId(),
			Operation_Id:      v.GetId(),
			ParentOperationId: v.GetParentOperationId(),
			Name:              v.GetName(),
			Date:              v.Date.AsTime(),
			Type:              int64(v.GetType()),
			Description:       v.GetDescription(),
			InstrumentUid:     v.GetInstrumentUid(),
			Figi:              v.GetFigi(),
			InstrumentType:    v.GetInstrumentType(),
			InstrumentKind:    v.GetInstrumentKind().String(),
			PositionUid:       v.GetPositionUid(),
			Payment:           ConvertPbToMoneyValue(v.GetPayment()),
			Price:             ConvertPbToMoneyValue(v.GetPrice()),
			Commission:        ConvertPbToMoneyValue(v.GetCommission()),
			Yield:             ConvertPbToMoneyValue(v.GetYield()),
			YieldRelative:     ConvertPbToQuatation(v.GetYieldRelative()),
			AccruedInt:        ConvertPbToMoneyValue(v.GetAccruedInt()),
			QuantityDone:      v.GetQuantityDone(),
			AssetUid:          v.GetAssetUid(),
		}
		transformOperations = append(transformOperations, transformOperation)
	}
	return transformOperations
}

func (c *AnalyticsServiceClient) GetAllAssetUids(client *investgo.Client) (map[string]string, error) {
	const op = "service.GetAllAssetUids"
	instrumentService := client.NewInstrumentsServiceClient()
	AssetsResponse, err := instrumentService.GetAssets()
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: could not get assets uid", op)
	}
	assetUidInstrumentUidMap := make(map[string]string)
	for _, v := range AssetsResponse.AssetsResponse.Assets {
		asset_uid := v.Uid

		for _, instrument := range v.Instruments {
			instrument_uid := instrument.Uid
			assetUidInstrumentUidMap[instrument_uid] = asset_uid
		}
	}
	return assetUidInstrumentUidMap, nil
}

func (c *InstrumentsServiceClient) GetFutureBy(client *investgo.Client, figi string) (Future, error) {
	if figi == "" {
		return Future{}, ErrEmptyFigi
	}
	instrumentService := client.NewInstrumentsServiceClient()
	futuresResponse, err := instrumentService.FutureByFigi(figi)
	if err != nil {
		return Future{}, e.WrapIfErr("can't get futures by figi", err)
	}
	resp := convertFuturePbToFuture(futuresResponse.FutureResponse.Instrument)
	return resp, nil
}

func convertFuturePbToFuture(futurePb *pb.Future) Future {
	return Future{
		Name:                    futurePb.GetName(),
		MinPriceIncrement:       ConvertPbToQuatation(futurePb.GetMinPriceIncrement()),
		MinPriceIncrementAmount: ConvertPbToQuatation(futurePb.GetMinPriceIncrementAmount()),
		AssetType:               futurePb.GetAssetType(),
		BasicAssetPositionUid:   futurePb.GetBasicAssetPositionUid(),
	}
}

func (c *InstrumentsServiceClient) GetBondByUid(client *investgo.Client, uid string) (Bond, error) {
	const op = "service.GetBondByUid"
	if uid == "" {
		return Bond{}, fmt.Errorf("op:%s, error: incorrect uid: can't be empty string", op)
	}
	instrumentService := client.NewInstrumentsServiceClient()
	bondResponse, err := instrumentService.BondByUid(uid)
	if err != nil {
		return Bond{}, fmt.Errorf("op:%s, error: can't get share by figi", op)
	}
	resp := convertBondPbToBond(bondResponse.BondResponse.Instrument)
	return resp, nil
}

func convertBondPbToBond(bondPb *pb.Bond) Bond {
	return Bond{
		AciValue: ConvertPbToMoneyValue(bondPb.GetAciValue()),
		Nominal:  ConvertPbToMoneyValue(bondPb.GetNominal()),
		Currency: bondPb.GetCurrency(),
	}
}

func (c *InstrumentsServiceClient) GetCurrencyBy(client *investgo.Client, figi string) (Currency, error) {
	const op = "service.GetCurrencyBy"
	if figi == "" {
		return Currency{}, ErrEmptyFigi
	}
	instrumentService := client.NewInstrumentsServiceClient()
	currencyResponse, err := instrumentService.CurrencyByFigi(figi)
	if err != nil {
		return Currency{}, fmt.Errorf("op:%s, error: can't get curency by figi", op)
	}

	resp := convertCurrencyPbToCurrency(currencyResponse.CurrencyResponse.Instrument)
	return resp, nil
}

func convertCurrencyPbToCurrency(currencyPb *pb.Currency) Currency {
	return Currency{
		Isin: currencyPb.GetIsoCurrencyName(),
	}
}

func (c *InstrumentsServiceClient) FindBy(client *investgo.Client, query string) ([]InstrumentShort, error) {
	const op = "service.FindBy"
	if query == "" {
		return nil, ErrEmptyQuery
	}
	instrumentServiceClient := client.NewInstrumentsServiceClient()
	findInstr, err := instrumentServiceClient.FindInstrument(query)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: could not find instrument by query", op)
	}
	resp := convertInstrumentShortPbToInstrumentShort(findInstr.FindInstrumentResponse.GetInstruments())
	return resp, nil
}

func convertInstrumentShortPbToInstrumentShort(instrumentShortPb []*pb.InstrumentShort) []InstrumentShort {
	instrumentShorts := make([]InstrumentShort, 0, len(instrumentShortPb))
	for _, instrShortPb := range instrumentShortPb {
		instrumentShorts = append(instrumentShorts, InstrumentShort{
			InstrumentType: instrShortPb.GetInstrumentType(),
			Uid:            instrShortPb.GetUid(),
			Figi:           instrShortPb.GetFigi(),
		})
	}
	return instrumentShorts
}

func (c *AnalyticsServiceClient) GetBondsActions(client *investgo.Client, instrumentUid string) (BondIdentIdentifiers, error) {
	const op = "service.GetBondActions"
	if instrumentUid == "" {
		return BondIdentIdentifiers{}, ErrEmptyInstrumentUid
	}
	var res BondIdentIdentifiers
	instrumentService := client.NewInstrumentsServiceClient()
	bondUid, err := instrumentService.BondByUid(instrumentUid)
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("op: %s, error: could not get bond actions", op)
	}
	res.Ticker = bondUid.BondResponse.Instrument.GetTicker()
	res.ClassCode = bondUid.BondResponse.Instrument.GetClassCode()
	res.Name = bondUid.BondResponse.Instrument.GetName()

	if bondUid.BondResponse.Instrument.GetBondType() == 1 {
		res.Replaced = true
	}
	res.Nominal = ConvertPbToMoneyValue(bondUid.BondResponse.Instrument.GetNominal())
	res.NominalCurrency = bondUid.Instrument.GetNominal().Currency
	return res, nil
}

func (c *AnalyticsServiceClient) GetLastPriceInPersentageToNominal(client *investgo.Client, instrumentUid string) (LastPriceResponse, error) {
	const op = "tinkoffApi.GetLastPriceInPercentageToNominal"
	if instrumentUid == "" {
		return LastPriceResponse{}, fmt.Errorf("%s: %w", op, ErrEmptyInstrumentUid)
	}
	marketDataClient := client.NewMarketDataServiceClient()
	lastPriceAnswer, err := marketDataClient.GetLastPrices([]string{instrumentUid})
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	// Проверка через LastPriceType нужна для ошибке при некоректном instrumentUid.
	// Т.к. тинькофф выдает нулевую ошибку при любом ошибочном instrumentUid
	if lastPriceAnswer.GetLastPricesResponse.LastPrices[0].Price.ToFloat() == 0 &&
		lastPriceAnswer.LastPrices[0].LastPriceType.Number() == 0 {
		return LastPriceResponse{}, fmt.Errorf("%s: haven't response for instrument %s", op, instrumentUid)
	}

	if len(lastPriceAnswer.LastPrices) == 0 {
		return LastPriceResponse{}, fmt.Errorf("%s: no price data for instrument %s", op, instrumentUid)
	}

	lastPrice := ConvertPbToQuatation(lastPriceAnswer.LastPrices[0].Price)
	resp := LastPriceResponse{
		LastPrice: lastPrice,
	}
	return resp, nil
}

func (c *InstrumentsServiceClient) GetShareCurrencyBy(client *investgo.Client, figi string) (ShareCurrencyByResponse, error) {
	const op = "service.GetShareCurrencyBy"
	if figi == "" {
		return ShareCurrencyByResponse{}, ErrEmptyFigi
	}
	instrumentService := client.NewInstrumentsServiceClient()
	shareResponse, err := instrumentService.ShareByFigi(figi)
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("op: %s, error: can't get share by figi", op)
	}
	var resp ShareCurrencyByResponse
	resp.Currency = shareResponse.ShareResponse.Instrument.GetCurrency()
	return resp, nil
}
