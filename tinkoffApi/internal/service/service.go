package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tinkoffApi/lib/e"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func (c *Client) GetAcc() (map[string]Account, error) {
	usersService := c.Client.NewUsersServiceClient()
	accounts := make(map[string]Account)
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL // ПОтом надо обработать закрытые счета(например ИИС)
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		return nil, errors.New("GetAcc: operationsService.GetOperationsByCursor" + err.Error())
	} else {
		accs := accsResp.GetAccounts()
		for _, acc := range accs {
			account := Account{
				Id:          acc.GetId(),
				Type:        acc.GetType().String(),
				Name:        acc.GetName(),
				OpenedDate:  acc.GetOpenedDate().AsTime(),
				ClosedDate:  acc.GetClosedDate().AsTime(),
				Status:      int64(acc.GetStatus()),
				AccessLevel: acc.GetAccessLevel().String(),
			}
			accounts[acc.GetId()] = account
		}
	}

	return accounts, nil
}

var ErrCloseAccount = errors.New("close account haven't portffolio positions")

type Client struct {
	ctx    context.Context
	Logg   investgo.Logger
	config *investgo.Config
	Client *investgo.Client
}

func New(ctx context.Context, logg investgo.Logger, config *investgo.Config) *Client {
	return &Client{
		ctx:    ctx,
		Logg:   logg,
		config: config,
	}
}

// TODO: FillCleint заполнялся в сервисе в основной программе. Тут придется переписывать каждую функию. Можно делать это в хэндлерах.
// Это будет аналогом аутентификации
// TODO: Протестировать возможность выводв ошибки неверного токена. Дать неверный токен и проверить ошибку от Тинькофф
func (c *Client) FillClient(token string) (err error) {
	const op = "sevrice.FillClient"
	defer func() { err = e.WrapIfErr(fmt.Sprintf("op:%s, description:can't create Tinkoff Client", op), err) }()

	c.config.Token = token

	if err = c.getClient(); err != nil {
		return err
	}
	return nil
}

func (c *Client) getClient() error {
	client, err := investgo.NewClient(c.ctx, *c.config, c.Logg)
	if err != nil {
		return e.WrapIfErr("can't connect with tinkoffApi client", err)
	}
	c.Client = client
	return nil
}

func (c *Client) GetPortf(request PortfolioRequest) (_ Portfolio, err error) {
	accountID := request.AccountID
	accountStatus := request.AccountStatus
	portfolio := Portfolio{}
	if accountStatus == 3 {
		return portfolio, ErrCloseAccount
	}
	operationsService := c.Client.NewOperationsServiceClient()
	id := accountID
	portfolioResp, err := operationsService.GetPortfolio(id,
		pb.PortfolioRequest_RUB)
	if err != nil {
		return portfolio, e.WrapIfErr("can't get portifolio positions from tinkoff Api", err)
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

func (c *Client) GetOperations(request OperationsRequest) (_ []Operation, err error) {
	defer func() { err = e.WrapIfErr("can't get opperations from tinkoffApi", err) }()
	accountID := request.AccountID
	date := request.Date
	resOpereaions := make([]*pb.OperationItem, 0)
	opereationsService := c.Client.NewOperationsServiceClient()
	operationsResp, err := opereationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
		AccountId: accountID,
		From:      date,
		To:        time.Now(),
		Limit:     1000,
	})
	if err != nil {
		return nil, err
	}
	operations := operationsResp.GetOperationsByCursorResponse.GetItems()
	resOpereaions = append(resOpereaions, operations...)
	nextCursor := operationsResp.NextCursor
	for nextCursor != "" {
		operationsResp, err := opereationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
			AccountId: accountID,
			Limit:     1000,
			Cursor:    nextCursor,
		})
		if err != nil {
			return nil, err
		} else {
			nextCursor = operationsResp.NextCursor
			operations := operationsResp.GetOperationsByCursorResponse.Items
			resOpereaions = append(resOpereaions, operations...)
		}
	}
	resp := c.convertOperationsPbToOperaions(resOpereaions)
	fmt.Printf("✓ Добавлено %v операции в Account.Operation по счету %s\n", len(resOpereaions), accountID)
	return resp, nil
}

func (c *Client) convertOperationsPbToOperaions(operations []*pb.OperationItem) []Operation {
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
			InstrumentKind:    string(v.GetInstrumentKind()),
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

func (c *Client) GetAllAssetUids() (map[string]string, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	AssetsResponse, err := instrumentService.GetAssets()
	if err != nil {
		return nil, errors.New("GetAllAssetUids: instrumentService.GetAssets" + err.Error())
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

func (c *Client) GetFutureBy(figi string) (Future, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
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

func (c *Client) GetBondByUid(uid string) (Bond, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	bondResponse, err := instrumentService.BondByUid(uid)
	if err != nil {
		return Bond{}, e.WrapIfErr("can't get share by figi", err)
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

// TODO : Проверить почему ничего не возвращает.
func (c *Client) GetCurrencyBy(figi string) (Currency, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	currencyResponse, err := instrumentService.CurrencyByFigi(figi)
	if err != nil {
		return Currency{}, e.WrapIfErr("can't get share by figi", err)
	}
	resp := convertCurrencyPbToCurrency(currencyResponse.CurrencyResponse.Instrument)
	return resp, nil
}

func convertCurrencyPbToCurrency(currencyPb *pb.Currency) Currency {
	return Currency{
		Isin: currencyPb.GetIsin(),
	}
}

func (c *Client) GetBaseShareFutureValute(positionUid string) (BaseShareFutureValuteResponse, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	instrumentsShortResponce, err := instrumentService.FindInstrument(positionUid)
	if err != nil {
		return BaseShareFutureValuteResponse{}, e.WrapIfErr("can't get base share future valute", err)
	}
	instrumentsShort := instrumentsShortResponce.Instruments
	if len(instrumentsShort) == 0 {
		return BaseShareFutureValuteResponse{}, errors.New("can't get base share future valute")
	}
	instrument := instrumentsShort[0]
	if instrument.InstrumentType != "share" {
		return BaseShareFutureValuteResponse{}, errors.New("instrument is not share")
	}
	currency, err := c.getShareCurrencyBy(instrument.Figi)
	if err != nil {
		return BaseShareFutureValuteResponse{}, e.WrapIfErr("can't get base share future valute", err)
	}

	var resp BaseShareFutureValuteResponse
	resp.Currency = currency.Currency
	return resp, nil
}

func (c *Client) FindBy(query string) ([]InstrumentShort, error) {
	client := c.Client.NewInstrumentsServiceClient()
	findInstr, err := client.FindInstrument(query)
	if err != nil {
		return nil, e.WrapIfErr("findByTicker error", err)
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
		})
	}
	return instrumentShorts
}

func (c *Client) GetBondsActions(instrumentUid string) (BondIdentIdentifiers, error) {
	var res BondIdentIdentifiers
	instrumentService := c.Client.NewInstrumentsServiceClient()
	bondUid, err := instrumentService.BondByUid(instrumentUid)
	if err != nil {
		return res, errors.New("GetTickerFromUid: instrumentService.BondByUid" + err.Error())
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

func (c *Client) GetLastPriceInPersentageToNominal(instrumentUid string) (LastPriceResponse, error) {
	marketDataClient := c.Client.NewMarketDataServiceClient()
	lastPriceAnswer, err := marketDataClient.GetLastPrices([]string{instrumentUid})
	if err != nil {
		return LastPriceResponse{}, errors.New("tinkoffApi:GetLastPriceFromTinkoff" + err.Error())
	}
	if len(lastPriceAnswer.LastPrices) == 0 {
		return LastPriceResponse{}, errors.New("tinkoffApi:GetLastPriceFromTinkoff: no price data")
	}

	lastPrice := ConvertPbToQuatation(lastPriceAnswer.LastPrices[0].Price)
	resp := LastPriceResponse{
		LastPrice: lastPrice,
	}
	return resp, nil
}

func (c *Client) getShareCurrencyBy(figi string) (ShareCurrencyByResponse, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	shareResponse, err := instrumentService.ShareByFigi(figi)
	if err != nil {
		return ShareCurrencyByResponse{}, e.WrapIfErr("can't get share by figi", err)
	}
	var resp ShareCurrencyByResponse
	resp.Currency = shareResponse.ShareResponse.Instrument.GetCurrency()
	return resp, nil
}
