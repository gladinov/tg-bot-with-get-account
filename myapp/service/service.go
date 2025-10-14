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
	"main.go/clients/sber"
	"main.go/clients/tinkoffApi"
	"main.go/lib/e"
	pathwd "main.go/lib/pathWD"
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

const (
	sberConfigPath = "/configs/sber.yaml"
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

func (c *Client) GetBondReports(chatID int, token string) (_ [][]*service_models.MediaGroup, err error) {
	defer func() { err = e.WrapIfErr("can't get general bond report", err) }()
	reportsInByteByAccounts := make([][]*service_models.MediaGroup, 0)
	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return nil, err
	}

	accounts, err := c.Tinkoffapi.GetAcc()
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		err = c.updateOperations(chatID, account.Id, account.OpenedDate)
		if err != nil {
			return nil, err
		}

		portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
		if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
			return nil, err
		}

		portfolioPositions, err := c.TransformPositions(account.Id, portfolio.Positions)
		if err != nil {
			return nil, err
		}
		err = c.Storage.DeleteGeneralBondReport(context.Background(), chatID, account.Id)
		if err != nil {
			return nil, err
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
					return nil, err
				}
				resultBondPosition, err := c.ProcessOperations(operationsDb)
				if err != nil {
					return nil, err
				}
				bondReport, err := c.CreateGeneralBondReport(resultBondPosition, portfolio.TotalAmount)
				if err != nil {
					return nil, err
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

		reportsInByte, err := c.PrepareToGenerateTablePNG(&generalBondReports, chatID, account.Id)
		if err != nil {
			return nil, err
		}
		reportsInByteByAccounts = append(reportsInByteByAccounts, reportsInByte)
		// err = c.Storage.SaveGeneralBondReport(context.Background(), chatID, account.Id, generalBondReportPositionsSorted)
		// if err != nil {
		// 	return err
		// }
	}

	return reportsInByteByAccounts, nil
}

func (c *Client) PrepareToGenerateTablePNG(generalBondReports *service_models.GeneralBondReports, chatID int, accountId string) (_ []*service_models.MediaGroup, err error) {
	reports := make([][]service_models.GeneralBondReportPosition, 0)

	rubbleBondReportSorted := sortGeneralBondReports(generalBondReports.RubBondsReport)
	replacedBondReportSorted := sortGeneralBondReports(generalBondReports.ReplacedBondsReport)
	euroBondReportSorted := sortGeneralBondReports(generalBondReports.EuroBondsReport)
	reports = append(reports, rubbleBondReportSorted)
	reports = append(reports, replacedBondReportSorted)
	reports = append(reports, euroBondReportSorted)
	reportsInByte := make([]*service_models.MediaGroup, 3)
	for i, report := range reports {
		reportsInByte[i] = service_models.NewMediaGroup()
		mediaGroup := reportsInByte[i]
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
		count := 1
		for start := 0; start < len(report); start += 10 {
			end := start + 10
			if end > len(report) {
				end = len(report)
			}
			pngData, err := visualization.GenerateTablePNG(report[start:end], typeOfBonds)
			if err != nil {
				return nil, e.WrapIfErr("vizualize error", err)
			}
			imageData := service_models.NewImageData()
			imageData.Name = fmt.Sprintf("file%s_%v", typeOfBonds, count)
			imageData.Data = pngData
			imageData.Caption = typeOfBonds

			mediaGroup.Reports = append(mediaGroup.Reports, imageData)
			count += 1
		}
	}
	return reportsInByte, nil
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
	if errors.Is(err, tinkoffApi.ErrCloseAccount) {
		return "", tinkoffApi.ErrCloseAccount
	}
	positions := portfolio.Positions

	accountTitle := fmt.Sprintf("Струтура брокерского счета: %s\n", account.Name)
	potfolioStructure, err := c.DivideByType(positions)
	if err != nil {
		return "", err
	}
	respPotfolioStructure := c.ResponsePortfolioStructure(potfolioStructure)

	response := accountTitle + respPotfolioStructure
	return response, nil
}

func (c *Client) GetUnionPortfolioStructure(token string, accounts map[string]tinkoffApi.Account) (_ string, err error) {
	defer func() { err = e.WrapIfErr("service: can't get union portfolio structure", err) }()
	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return "", err
	}

	positionsList := make([]*service_models.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
		if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
			return "", err
		}
		if errors.Is(err, tinkoffApi.ErrCloseAccount) {
			continue
		}
		positions := portfolio.Positions

		potfolioStructure, err := c.DivideByType(positions)
		if err != nil {
			return "", err
		}
		positionsList = append(positionsList, potfolioStructure)
	}
	accountTitle := "Струтура всех открытых счетов в Тинькофф Инвестициях\n"
	unionPositions, err := c.UnionPortf(positionsList)
	if err != nil {
		return "", err
	}
	vizualizeUnionPositions := c.ResponsePortfolioStructure(unionPositions)
	out := accountTitle + vizualizeUnionPositions
	return out, nil
}

func (c *Client) GetUnionPortfolioStructureWithSber(token string, accounts map[string]tinkoffApi.Account) (_ string, err error) {
	defer func() { err = e.WrapIfErr("service: can't get union portfolio structure with Sber", err) }()
	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return "", err
	}

	positionsList := make([]*service_models.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
		if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
			return "", err
		}
		if errors.Is(err, tinkoffApi.ErrCloseAccount) {
			continue
		}
		positions := portfolio.Positions

		potfolioStructure, err := c.DivideByType(positions)
		if err != nil {
			return "", err
		}
		positionsList = append(positionsList, potfolioStructure)
	}

	sberConfigAbsolutPath, err := pathwd.PathFromWD(sberConfigPath)
	if err != nil {
		return "", err
	}
	sberConfig, err := sber.LoadConfigSber(sberConfigAbsolutPath)
	if err != nil {
		return "", err
	}

	processConfig, err := sber.ProcessConfigSber(sberConfig)
	if err != nil {
		return "", err
	}

	sberPortfolio, err := c.DivideByTypeFromSber(processConfig)
	if err != nil {
		return "", err
	}

	positionsList = append(positionsList, sberPortfolio)

	accountTitle := "Струтура всех инвестиций\n"
	unionPositions, err := c.UnionPortf(positionsList)
	if err != nil {
		return "", err
	}
	vizualizeUnionPositions := c.ResponsePortfolioStructure(unionPositions)
	out := accountTitle + vizualizeUnionPositions
	return out, nil
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

func (c *Client) DivideByTypeFromSber(positions map[string]float64) (*service_models.PortfolioByTypeAndCurrency, error) {
	portfolio := service_models.NewPortfolioByTypeAndCurrency()

	if len(positions) == 0 {
		return portfolio, errors.New("positions are empty")
	}
	for ticker, quantity := range positions {
		positionsClassCodeVariants, err := c.Tinkoffapi.FindBy(ticker)
		if err != nil {
			return nil, e.WrapIfErr("can't divide by type from sber", err)
		}
		if len(positionsClassCodeVariants) == 0 {
			return nil, errors.New("positions variants are empty")
		}

		switch positionsClassCodeVariants[0].InstrumentType {
		case bond:
			bondUid := positionsClassCodeVariants[0].Uid
			bondActions, err := c.Tinkoffapi.GetBondByUid(bondUid)
			if err != nil {
				return nil, e.WrapIfErr("can't divide by type from sber", err)
			}
			currentNkd := bondActions.AciValue.ToFloat()
			currency := bondActions.Currency
			currentPriceInPers, err := c.Tinkoffapi.GetLastPriceFromTinkoffInPersentageToNominal(bondUid)
			if err != nil {
				return nil, e.WrapIfErr("can't divide by type from sber", err)
			}
			currentPrice := currentPriceInPers / 100 * bondActions.Nominal.ToFloat()
			currentNkdOfPosition := currentNkd * quantity
			positionPrice := currentPrice*quantity + currentNkdOfPosition

			portfolio.AllAssets += positionPrice
			portfolio.BondsAssets.SumOfAssets += positionPrice

			if existing, exist := portfolio.BondsAssets.AssetsByCurrency[currency]; !exist {
				portfolio.BondsAssets.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: positionPrice,
				}
			} else {
				existing.SumOfAssets += positionPrice
			}

		case share:
		case futures:
		case etf:
		case currency:
		default:
		}

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

func (c *Client) UnionPortf(portfolios []*service_models.PortfolioByTypeAndCurrency) (*service_models.PortfolioByTypeAndCurrency, error) {
	unionPortf := service_models.NewPortfolioByTypeAndCurrency()
	for _, portf := range portfolios {
		unionPortf.AllAssets += portf.AllAssets

		unionPortf.BondsAssets.SumOfAssets += portf.BondsAssets.SumOfAssets
		for k, v := range portf.BondsAssets.AssetsByCurrency {
			if existing, exist := unionPortf.BondsAssets.AssetsByCurrency[k]; !exist {
				unionPortf.BondsAssets.AssetsByCurrency[k] = service_models.NewAssetsByParam()
				unionPortf.BondsAssets.AssetsByCurrency[k] = v
			} else {
				existing.SumOfAssets += v.SumOfAssets
			}
		}

		unionPortf.SharesAssets.SumOfAssets += portf.SharesAssets.SumOfAssets
		for currency, asset := range portf.SharesAssets.AssetsByCurrency {
			if existing, exists := unionPortf.SharesAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.SharesAssets.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.EtfsAssets.SumOfAssets += portf.EtfsAssets.SumOfAssets
		for currency, asset := range portf.EtfsAssets.AssetsByCurrency {
			if existing, exists := unionPortf.EtfsAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.EtfsAssets.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.CurrenciesAssets.SumOfAssets += portf.CurrenciesAssets.SumOfAssets
		for currency, asset := range portf.CurrenciesAssets.AssetsByCurrency {
			if existing, exists := unionPortf.CurrenciesAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.CurrenciesAssets.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.SumOfAssets += portf.FuturesAssets.SumOfAssets
		unionPortf.FuturesAssets.AssetsByType.Commodity.SumOfAssets += portf.FuturesAssets.AssetsByType.Commodity.SumOfAssets
		for currency, asset := range portf.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency {
			if existing, exist := unionPortf.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[currency]; !exist {
				unionPortf.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Currency.SumOfAssets += portf.FuturesAssets.AssetsByType.Currency.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Security.SumOfAssets += portf.FuturesAssets.AssetsByType.Security.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Security.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Security.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Security.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Index.SumOfAssets += portf.FuturesAssets.AssetsByType.Index.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Index.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Index.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Index.AssetsByCurrency[currency] = &service_models.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}
	}
	return unionPortf, nil
}
