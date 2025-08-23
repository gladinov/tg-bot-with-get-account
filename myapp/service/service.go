package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"main.go/clients/cbr"
	"main.go/clients/moex"
	"main.go/clients/tinkoffApi"
	"main.go/lib/e"
	"main.go/service/service_models"
	service_storage "main.go/service/storage"
	"main.go/service/visualization"
)

const (
	layoutTime = "2006-01-02_15-04-05"
	layoutCurr = "02.01.2006"
	reportPath = "service/visualization/tables/"
)

const (
	bond     = "bond"
	share    = "share"
	futures  = "futures"
	etf      = "etf"
	currency = "currency"
)

const (
	rub       = "rub"
	cny       = "cny"
	usd       = "usd"
	eur       = "eur"
	hkd       = "hkd"
	futuresPt = "pt."
)

const (
	commodityType = "TYPE_COMMODITY"
	currencyType  = "TYPE_CURRENCY"
	securityType  = "TYPE_SECURITY"
	indexType     = "TYPE_INDEX"
)

type Client struct {
	Tinkoffapi *tinkoffApi.Client
	MoexApi    *moex.Client
	CbrApi     *cbr.Client
	Storage    service_storage.Storage
}

func New(tinkoffApiClient *tinkoffApi.Client, moexClient *moex.Client, CbrClient *cbr.Client, storage service_storage.Storage) *Client {
	return &Client{
		Tinkoffapi: tinkoffApiClient,
		MoexApi:    moexClient,
		CbrApi:     CbrClient,
		Storage:    storage,
	}
}

func (c *Client) GetBondReportsByFifo(chatID int, token string) (err error) {
	defer func() { err = e.WrapIfErr("can't get bond reports", err) }()

	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return err
	}

	accounts, err := c.Tinkoffapi.GetAcc()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = c.updateOperations(chatID, account.Id, account.OpenedDate)
		if err != nil {
			return err
		}

		portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
		if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
			return err
		}

		portfolioPositions, err := c.TransformPositions(account.Id, portfolio.Positions)
		if err != nil {
			return err
		}
		err = c.Storage.DeleteBondReport(context.Background(), chatID, account.Id)
		if err != nil {
			return err
		}
		bondsInRub := make([]service_models.BondReport, 0)
		for _, v := range portfolioPositions {
			if v.InstrumentType == "bond" {
				operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.AssetUid, account.Id)
				if err != nil {
					return err
				}
				resultBondPosition, err := c.ProcessOperations(operationsDb)
				if err != nil {
					return err
				}
				bondReport, err := c.CreateBondReport(*resultBondPosition)
				if err != nil {
					return err
				}
				bondsInRub = append(bondsInRub, bondReport.BondsInRUB...)
			}
		}
		err = c.Storage.SaveBondReport(context.Background(), chatID, account.Id, bondsInRub)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GetBondReportsWithEachGeneralPosition(chatID int, token string) (err error) {
	defer func() { err = e.WrapIfErr("can't get general bond report", err) }()
	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return err
	}

	accounts, err := c.Tinkoffapi.GetAcc()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = c.updateOperations(chatID, account.Id, account.OpenedDate)
		if err != nil {
			return err
		}

		portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
		if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
			return err
		}

		portfolioPositions, err := c.TransformPositions(account.Id, portfolio.Positions)
		if err != nil {
			return err
		}
		err = c.Storage.DeleteGeneralBondReport(context.Background(), chatID, account.Id)
		if err != nil {
			return err
		}
		generalBondReports := service_models.GeneralBondReports{
			RubBondsReport:      make(map[service_models.TickerTimeKey]service_models.GeneralBondReportPosition),
			EuroBondsReport:     make(map[service_models.TickerTimeKey]service_models.GeneralBondReportPosition),
			ReplacedBondsReport: make(map[service_models.TickerTimeKey]service_models.GeneralBondReportPosition),
		}

		for _, v := range portfolioPositions {
			if v.InstrumentType == "bond" {
				operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.AssetUid, account.Id)
				if err != nil {
					return err
				}
				resultBondPosition, err := c.ProcessOperations(operationsDb)
				if err != nil {
					return err
				}
				bondReport, err := c.CreateGeneralBondReport(resultBondPosition, portfolio.TotalAmount)
				if err != nil {
					return err
				}
				switch {
				case bondReport.Replaced:
					tickerTimeKey := service_models.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.ReplacedBondsReport[tickerTimeKey] = bondReport
				case bondReport.Currencies != "rub":
					tickerTimeKey := service_models.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.EuroBondsReport[tickerTimeKey] = bondReport
				default:
					tickerTimeKey := service_models.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.RubBondsReport[tickerTimeKey] = bondReport
				}

			}
		}

		err = Vizualization(&generalBondReports, chatID, account.Id)
		if err != nil {
			return err
		}

		// err = c.Storage.SaveGeneralBondReport(context.Background(), chatID, account.Id, generalBondReportPositionsSorted)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func Vizualization(generalBondReports *service_models.GeneralBondReports, chatID int, accountId string) error {
	reports := make([][]service_models.GeneralBondReportPosition, 0)

	rubbleBondReportSorted := sortGeneralBondReports(generalBondReports.RubBondsReport)
	replacedBondReportSorted := sortGeneralBondReports(generalBondReports.ReplacedBondsReport)
	euroBondReportSorted := sortGeneralBondReports(generalBondReports.EuroBondsReport)
	reports = append(reports, rubbleBondReportSorted)
	reports = append(reports, replacedBondReportSorted)
	reports = append(reports, euroBondReportSorted)

	for _, report := range reports {
		if len(report) == 0 {
			continue
		}

		var typeOfBonds string
		switch {
		case report[0].Replaced:
			typeOfBonds = service_models.ReplacedBonds
		case report[0].Currencies != "rub":
			typeOfBonds = service_models.EuroBonds
		default:
			typeOfBonds = service_models.RubBonds
		}
		pathDir := path.Join(reportPath, strconv.Itoa(chatID), accountId)
		if _, err := os.Stat(pathDir); os.IsNotExist(err) {
			err = os.MkdirAll(pathDir, 0755)
			if err != nil {
				return e.WrapIfErr("can't make directory", err)
			}
		}
		count := 1
		now := time.Now()
		nameTime := now.Format(layoutTime)

		for start := 0; start < len(report); start += 10 {
			countName := strconv.Itoa(count)
			fileName := nameTime + "_" + typeOfBonds + "_" + countName + ".png"
			pathAndName := path.Join(pathDir, fileName)
			end := start + 10
			if end > len(report) {
				end = len(report)
			}
			err := visualization.Vizualize(report[start:end], pathAndName, typeOfBonds)
			if err != nil {
				return e.WrapIfErr("vizualize error", err)
			}
			count += 1

		}
	}
	return nil
}

func sortGeneralBondReports(report map[service_models.TickerTimeKey]service_models.GeneralBondReportPosition) []service_models.GeneralBondReportPosition {
	keys := make([]service_models.TickerTimeKey, 0, len(report))
	for k := range report {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Time.Equal(keys[j].Time) {
			return keys[i].Ticker < keys[j].Ticker
		}
		return keys[i].Time.Before(keys[j].Time)
	})
	result := make([]service_models.GeneralBondReportPosition, len(keys))
	for i, k := range keys {
		result[i] = report[k]
	}

	return result
}

func (c *Client) GetAccountsList(token string) (answ string, err error) {
	defer func() { err = e.WrapIfErr("can't get accounts", err) }()
	client := c.Tinkoffapi
	var accStr string = "По данному аккаунту доступны следующие счета:"
	err = client.FillClient(token)
	if err != nil {
		return "", err
	}

	accs, err := c.Tinkoffapi.GetAcc()
	if err != nil {
		return "", err
	}
	for _, account := range accs {
		accStr += fmt.Sprintf("\n ID:%v, Type: %s, Name: %s, Status: %v \n", account.Id, account.Type, account.Name, account.Status)
	}

	return accStr, nil
}

func (c *Client) GetUsd() (float64, error) {

	usd, err := c.GetCurrencyFromCB("usd", time.Now())
	if err != nil {
		return 0, e.WrapIfErr("usd get error", err)
	}
	return usd, nil
}

func (c *Client) updateOperations(chatID int, accountId string, openDate time.Time) (err error) {
	defer func() { err = e.WrapIfErr("can't updateOperations", err) }()
	fromDate, err := c.Storage.LastOperationTime(context.Background(), chatID, accountId)
	fromDate = fromDate.Add(time.Microsecond * 1)

	if err != nil {
		if errors.Is(err, service_models.ErrNoOpperations) {
			fromDate = openDate
		} else {
			return err
		}
	}

	tinkoffOperations, err := c.Tinkoffapi.GetOperations(accountId, fromDate)
	if err != nil {
		return err
	}
	operations := c.TransOperations(tinkoffOperations)

	err = c.Storage.SaveOperations(context.Background(), chatID, accountId, operations)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccounts(token string) (_ map[string]tinkoffApi.Account, err error) {
	defer func() { err = e.WrapIfErr("cant' get accounts", err) }()
	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return nil, err
	}

	accounts, err := c.Tinkoffapi.GetAcc()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (c *Client) GetPortfolioStructure(token string, account tinkoffApi.Account) (_ string, err error) {
	defer func() { err = e.WrapIfErr("cant' get portfolio structure", err) }()
	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return "", err
	}

	portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
	if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
		return "", err
	}
	positions := portfolio.Positions
	accountTitle := fmt.Sprintf("Струтура брокерского счета: %s\n", account.Name)
	potfolioStructure, err := c.DivideByType(positions)
	if err != nil {
		return "", nil
	}
	respPotfolioStructure := c.ResponsePortfolioStructure(potfolioStructure)

	response := accountTitle + respPotfolioStructure
	return response, nil
}

func (c *Client) DivideByType(positions []*pb.PortfolioPosition) (_ *service_models.PortfolioByTypeAndCurrency, err error) {
	defer func() { err = e.WrapIfErr("can't divide by type", err) }()
	portfolio := service_models.NewPortfolioByTypeAndCurrency()
	date := time.Now()

	if len(positions) == 0 {
		return portfolio, errors.New("positions are empty")
	}

	for _, pos := range positions {
		var positionPrice float64
		currencyOfPos := pos.CurrentPrice.Currency
		vunit_rate := 1.0
		if currencyOfPos != futuresPt && currencyOfPos != rub {
			vunit_rate, err = c.GetCurrencyFromCB(currencyOfPos, date)
			if err != nil {
				return portfolio, e.WrapIfErr("can't divide by type", err)
			}
		}
		positionPrice = pos.Quantity.ToFloat() * pos.CurrentPrice.ToFloat() * vunit_rate

		switch pos.InstrumentType {
		case bond:
			positionPrice += pos.CurrentNkd.ToFloat() * pos.GetQuantity().ToFloat() * vunit_rate
			portfolio.BondsAssets.SumOfAssets += positionPrice
			if _, exist := portfolio.BondsAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.BondsAssets.AssetsByCurrency[currencyOfPos] = service_models.NewAssetsByParam()
			}
			portfolio.BondsAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case share:
			portfolio.SharesAssets.SumOfAssets += positionPrice
			if _, exist := portfolio.SharesAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.SharesAssets.AssetsByCurrency[currencyOfPos] = service_models.NewAssetsByParam()
			}
			portfolio.SharesAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case futures:

			futures, err := c.Tinkoffapi.GetFutureBy(pos.Figi)
			if err != nil {
				return portfolio, e.WrapIfErr("can't divide by type", err)
			}

			positionPrice = positionPrice / futures.MinPriceIncrement.ToFloat() * futures.MinPriceIncrementAmount.ToFloat()
			portfolio.FuturesAssets.SumOfAssets += positionPrice

			futureType := futures.AssetType
			switch futureType {
			case commodityType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name] = service_models.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Commodity.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice
			case currencyType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name] = service_models.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Currency.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice
			case securityType:
				valute, err := c.Tinkoffapi.GetBaseShareFutureValute(futures.BasicAssetPositionUid)
				if err != nil {
					return nil, e.WrapIfErr("can't divide by type", err)
				}
				if _, exist := portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute]; !exist {
					portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute] = service_models.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Security.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute].SumOfAssets += positionPrice
			case indexType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name] = service_models.NewAssetsByParam()
				}

				portfolio.FuturesAssets.AssetsByType.Index.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice

			}
			// Чтобы сумма фьюча не сумировалась с суммой всех активов, так как фактически я за тело фьючерса не заплатил
			positionPrice = 0

		case etf:
			portfolio.EtfsAssets.SumOfAssets += positionPrice

			if _, exist := portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos] = service_models.NewAssetsByParam()
			}
			portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case currency:
			curr, err := c.Tinkoffapi.GetCurrencyBy(pos.Figi)
			if err != nil {
				return portfolio, e.WrapIfErr("can't divide by type", err)
			}
			currName := curr.Isin
			portfolio.CurrenciesAssets.SumOfAssets += positionPrice

			if _, exist := portfolio.CurrenciesAssets.AssetsByCurrency[currName]; !exist {
				portfolio.CurrenciesAssets.AssetsByCurrency[currName] = service_models.NewAssetsByParam()
			}
			portfolio.CurrenciesAssets.AssetsByCurrency[currName].SumOfAssets += positionPrice

		default:
		}
		portfolio.AllAssets += positionPrice
	}

	return portfolio, nil
}

func (c *Client) ResponsePortfolioStructure(portfolio *service_models.PortfolioByTypeAndCurrency) string {
	var output string
	totalAmount := portfolio.AllAssets
	totalBonds := portfolio.BondsAssets.SumOfAssets
	totalShares := portfolio.SharesAssets.SumOfAssets
	totalEtfs := portfolio.EtfsAssets.SumOfAssets
	totalFutures := portfolio.FuturesAssets.SumOfAssets
	totalCurrencies := portfolio.CurrenciesAssets.SumOfAssets
	totalLeverageRatio := (totalBonds + totalShares + totalEtfs + totalFutures) / totalAmount

	totalBondsInPerc := totalBonds / totalAmount * 100
	totalSharesInPerc := totalShares / totalAmount * 100
	totalEtfsInPerc := totalEtfs / totalAmount * 100
	totalFuturesInPerc := totalFutures / totalAmount * 100
	totalCurrenciesInPerc := totalCurrencies / totalAmount * 100
	totalLeverageRatioInPers := totalLeverageRatio * 100

	totalAmountOut := fmt.Sprintf("Общая стоимость портфеля составляет: %.2f\n", totalAmount)
	totalBondsOut := fmt.Sprintf("  Стоимость облигаций: %.2f (%.2f%%от портфеля)\n", totalBonds, totalBondsInPerc)
	totalSharesOut := fmt.Sprintf("  Стоимость акций: %.2f (%.2f%%от портфеля)\n", totalShares, totalSharesInPerc)
	totalEtfsOut := fmt.Sprintf("  Стоимость ETF: %.2f (%.2f%%от портфеля)\n", totalEtfs, totalEtfsInPerc)
	totalFuturesOut := fmt.Sprintf("  Стоимость фьючерсов: %.2f (%.2f%%от портфеля)\n", totalFutures, totalFuturesInPerc)
	totalCurrenciesOut := fmt.Sprintf("  Стоимость валют: %.2f (%.2f%%от портфеля)\n", totalCurrencies, totalCurrenciesInPerc)
	totalLeverageRatioOut := fmt.Sprintf("  Общий коэффициент левериджа: %.2f (%.2f%%от портфеля)\n", totalLeverageRatio, totalLeverageRatioInPers)

	var bondsByCurrencies string
	for k, v := range portfolio.BondsAssets.AssetsByCurrency {
		bondsByCurr := RoundFloat(v.SumOfAssets, 2)
		bondsByCurrInPers := RoundFloat(bondsByCurr/totalAmount*100, 2)
		bondsByCurrencies += "    " + fmt.Sprintf("Стоимость облигаций в %s: %.2f (%.2f%%от портфеля)\n", k, bondsByCurr, bondsByCurrInPers)
	}

	var sharesByCurrencies string
	for k, v := range portfolio.SharesAssets.AssetsByCurrency {
		AssetByCurr := RoundFloat(v.SumOfAssets, 2)
		AssetByCurrInPers := RoundFloat(AssetByCurr/totalAmount*100, 2)
		sharesByCurrencies += "    " + fmt.Sprintf("Стоимость акций в %s: %.2f (%.2f%%от портфеля)\n", k, AssetByCurr, AssetByCurrInPers)
	}

	var etfsByCurrencies string
	for k, v := range portfolio.EtfsAssets.AssetsByCurrency {
		AssetByCurr := RoundFloat(v.SumOfAssets, 2)
		AssetByCurrInPers := RoundFloat(AssetByCurr/totalAmount*100, 2)
		etfsByCurrencies += "    " + fmt.Sprintf("Стоимость Etf в %s: %.2f (%.2f%%от портфеля)\n", k, AssetByCurr, AssetByCurrInPers)
	}

	var futuresByCurrencies string

	if portfolio.FuturesAssets.AssetsByType.Commodity.SumOfAssets != 0 {
		futuresWithCommodityBase := portfolio.FuturesAssets.AssetsByType.Commodity.SumOfAssets
		futuresWithCommodityBaseInPers := futuresWithCommodityBase / totalAmount * 100

		futuresByCurrencies += "    " + fmt.Sprintf("Фьючерсы на товары стоят: %.2f (%.2f%%от портфеля)\n", futuresWithCommodityBase, futuresWithCommodityBaseInPers)
		for k, v := range portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency {
			AssetByCurr := RoundFloat(v.SumOfAssets, 2)
			AssetByCurrInPers := RoundFloat(AssetByCurr/totalAmount*100, 2)
			futuresByCurrencies += "      " + fmt.Sprintf("Стоимость фьючерса %s: %.2f (%.2f%%от портфеля)\n", k, AssetByCurr, AssetByCurrInPers)
		}
	}

	if portfolio.FuturesAssets.AssetsByType.Currency.SumOfAssets != 0 {
		futuresWithcurrencyBase := portfolio.FuturesAssets.AssetsByType.Currency.SumOfAssets
		futuresWithcurrencyBaseInPers := futuresWithcurrencyBase / totalAmount * 100

		futuresByCurrencies += "    " + fmt.Sprintf("Фьючерсы на валюты стоят: %.2f (%.2f%%от портфеля)\n", futuresWithcurrencyBase, futuresWithcurrencyBaseInPers)

		for k, v := range portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency {
			AssetByCurr := RoundFloat(v.SumOfAssets, 2)
			AssetByCurrInPers := RoundFloat(AssetByCurr/totalAmount*100, 2)
			futuresByCurrencies += "      " + fmt.Sprintf("Стоимость фьючерса %s: %.2f (%.2f%%от портфеля)\n", k, AssetByCurr, AssetByCurrInPers)
		}

	}
	if portfolio.FuturesAssets.AssetsByType.Security.SumOfAssets != 0 {
		futuresWithcurrencyBase := portfolio.FuturesAssets.AssetsByType.Security.SumOfAssets
		futuresWithcurrencyBaseInPers := futuresWithcurrencyBase / totalAmount * 100

		futuresByCurrencies += "    " + fmt.Sprintf("Фьючерсы на акции стоят: %.2f (%.2f%%от портфеля)\n", futuresWithcurrencyBase, futuresWithcurrencyBaseInPers)

		for k, v := range portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency {
			AssetByCurr := RoundFloat(v.SumOfAssets, 2)
			AssetByCurrInPers := RoundFloat(AssetByCurr/totalAmount*100, 2)
			futuresByCurrencies += "      " + fmt.Sprintf("Стоимость фьючерсов в %s: %.2f (%.2f%%от портфеля)\n", k, AssetByCurr, AssetByCurrInPers)
		}
	}

	if portfolio.FuturesAssets.AssetsByType.Index.SumOfAssets != 0 {
		futuresWithcurrencyBase := portfolio.FuturesAssets.AssetsByType.Index.SumOfAssets
		futuresWithcurrencyBaseInPers := futuresWithcurrencyBase / totalAmount * 100

		futuresByCurrencies += "    " + fmt.Sprintf("Фьючерсы на индексы стоят: %.2f (%.2f%%от портфеля)\n", futuresWithcurrencyBase, futuresWithcurrencyBaseInPers)
		futuresByCurrencies += "    " + "Фьючерсы на индексы\n"
		for k, v := range portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency {
			AssetByCurr := RoundFloat(v.SumOfAssets, 2)
			AssetByCurrInPers := RoundFloat(AssetByCurr/totalAmount*100, 2)
			futuresByCurrencies += "      " + fmt.Sprintf("Стоимость фьючерса %s: %.2f (%.2f%%от портфеля)\n", k, AssetByCurr, AssetByCurrInPers)
		}
	}
	var currenciesByCurrencies string
	for k, v := range portfolio.CurrenciesAssets.AssetsByCurrency {
		AssetByCurr := RoundFloat(v.SumOfAssets, 2)
		AssetByCurrInPers := RoundFloat(AssetByCurr/totalAmount*100, 2)
		currenciesByCurrencies += "      " + fmt.Sprintf("Стоимость валюты %s: %.2f (%.2f%%от портфеля)\n", k, AssetByCurr, AssetByCurrInPers)
	}

	output += totalAmountOut +
		totalBondsOut +
		bondsByCurrencies +
		totalSharesOut +
		sharesByCurrencies +
		totalEtfsOut +
		etfsByCurrencies +
		totalFuturesOut +
		futuresByCurrencies +
		totalCurrenciesOut +
		currenciesByCurrencies +
		totalLeverageRatioOut

	return output
}
