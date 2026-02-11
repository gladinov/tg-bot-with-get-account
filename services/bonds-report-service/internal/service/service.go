package service

import (
	"bonds-report-service/internal/clients/sber"
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/models/domain/mapper"
	"bonds-report-service/internal/service/visualization"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/gladinov/e"

	service_storage "bonds-report-service/internal/repository"
)

const (
	layoutTime = "2006-01-02_15-04-05"

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

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=UidProvider
type UidProvider interface {
	GetUid(ctx context.Context, instrumentUid string) (string, error)
	UpdateAndGetUid(ctx context.Context, instrumentUid string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=CbrClient
type CbrClient interface {
	GetAllCurrencies(ctx context.Context, date time.Time) (res domain.CurrenciesCBR, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=MoexClient
type MoexClient interface {
	GetSpecifications(ctx context.Context, ticker string, date time.Time) (data domain.ValuesMoex, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TinkoffInstrumentsClient
type TinkoffInstrumentsClient interface {
	FindBy(ctx context.Context, query string) ([]domain.InstrumentShort, error)
	GetBondByUid(ctx context.Context, uid string) (domain.Bond, error)
	GetCurrencyBy(ctx context.Context, figi string) (domain.Currency, error)
	GetFutureBy(ctx context.Context, figi string) (domain.Future, error)
	GetShareCurrencyBy(ctx context.Context, figi string) (domain.ShareCurrency, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TinkoffPortfolioClient
type TinkoffPortfolioClient interface {
	GetAccounts(ctx context.Context) (_ map[string]domain.Account, err error)
	GetPortfolio(ctx context.Context, accountID string, accountStatus int64) (domain.Portfolio, error)
	GetOperations(ctx context.Context, accountId string, date time.Time) (_ []domain.Operation, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TinkoffAnalyticsClient
type TinkoffAnalyticsClient interface {
	GetLastPriceInPersentageToNominal(ctx context.Context, instrumentUid string) (domain.LastPrice, error)
	GetAllAssetUids(ctx context.Context) (map[string]string, error)
	GetBondsActions(ctx context.Context, instrumentUid string) (domain.BondIdentIdentifiers, error)
}

type TinkoffClients struct {
	Instruments TinkoffInstrumentsClient
	Portfolio   TinkoffPortfolioClient
	Analytics   TinkoffAnalyticsClient
}

func NewTinkoffClients(
	instruments TinkoffInstrumentsClient,
	portfolio TinkoffPortfolioClient,
	analytics TinkoffAnalyticsClient,
) *TinkoffClients {
	return &TinkoffClients{
		Instruments: instruments,
		Portfolio:   portfolio,
		Analytics:   analytics,
	}
}

type ExternalApis struct {
	Moex MoexClient
	Cbr  CbrClient
	Sber *sber.Client
}

func NewExternalApis(
	moex MoexClient,
	cbr CbrClient,
	sber *sber.Client,
) *ExternalApis {
	return &ExternalApis{
		Moex: moex,
		Cbr:  cbr,
		Sber: sber,
	}
}

type Client struct {
	logger      *slog.Logger
	Tinkoff     *TinkoffClients
	External    *ExternalApis
	Storage     service_storage.Storage
	UidProvider UidProvider
}

func NewClient(
	logger *slog.Logger,
	tinkoffClients *TinkoffClients,
	externalApis *ExternalApis,
	storage service_storage.Storage,
	uidProvider UidProvider,
) *Client {
	return &Client{
		logger:      logger,
		Tinkoff:     tinkoffClients,
		External:    externalApis,
		Storage:     storage,
		UidProvider: uidProvider,
	}
}

func (c *Client) GetBondReportsByFifo(ctx context.Context, chatID int) (err error) {
	const op = "service.GetBondReportsByFifo"
	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get bond reports", err)
	}()
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		logg.Debug("get accounts error", slog.Any("error", err))
		return err
	}

	for _, account := range accounts {
		accountLogg := logg.With(
			slog.String("account_id", account.ID),
			slog.String("account_name", account.Name),
			slog.String("account_type", account.Type))
		err = c.updateOperations(ctx, chatID, account.ID, account.OpenedDate)
		if err != nil {
			accountLogg.Debug(
				"update operation error", slog.Any("error", err))
			return err
		}
		if account.Status != 2 {
			continue
		}
		portfolio, err := c.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			accountLogg.Debug(
				"tinkoffGetPortfolio err", slog.Any("error", err))
			return err
		}

		portfolioPositions, err := c.MapPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
		if err != nil {
			accountLogg.Debug(
				"transformPositions err", slog.Any("error", err))
			return err
		}
		err = c.Storage.DeleteBondReport(context.Background(), chatID, account.ID)
		if err != nil {
			accountLogg.Debug(
				"deleteBondReport err", slog.Any("error", err))
			return err
		}
		bondsInRub := make([]domain.BondReport, 0)

		for _, v := range portfolioPositions {
			positionLogg := accountLogg.With(
				slog.String("Asset_uid", v.AssetUid),
				slog.String("Instrument_type", v.InstrumentType))
			if v.InstrumentType == "bond" {
				operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.AssetUid, account.ID)
				if err != nil {
					positionLogg.Debug(
						"storage.GetOperation error", slog.Any("error", err))
					return err
				}
				resultBondPosition, err := c.ProcessOperations(ctx, operationsDb)
				if err != nil {
					positionLogg.Debug(
						"ProcessOperation error", slog.Any("error", err))
					return err
				}
				bondReport, err := c.CreateBondReport(ctx, *resultBondPosition)
				if err != nil {
					positionLogg.Debug(
						"CreateBondReport error", slog.Any("error", err))
					return err
				}
				bondsInRub = append(bondsInRub, bondReport.BondsInRUB...)
			}
		}
		err = c.Storage.SaveBondReport(context.Background(), chatID, account.ID, bondsInRub)
		if err != nil {
			accountLogg.Debug(
				"Storage.SaveBondReport error", slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (c *Client) GetBondReportsWithEachGeneralPosition(ctx context.Context, chatID int) (err error) {
	const op = "service.GetBondReportsWithEachGeneralPosition"
	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get general bond report", err)
	}()

	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = c.updateOperations(ctx, chatID, account.ID, account.OpenedDate)
		if err != nil {
			return err
		}
		if account.Status != 2 {
			continue
		}
		portfolio, err := c.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return err
		}

		portfolioPositions, err := c.MapPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
		if err != nil {
			return err
		}
		err = c.Storage.DeleteGeneralBondReport(context.Background(), chatID, account.ID)
		if err != nil {
			return err
		}
		generalBondReports := domain.GeneralBondReports{
			RubBondsReport:      make(map[domain.TickerTimeKey]domain.GeneralBondReportPosition),
			EuroBondsReport:     make(map[domain.TickerTimeKey]domain.GeneralBondReportPosition),
			ReplacedBondsReport: make(map[domain.TickerTimeKey]domain.GeneralBondReportPosition),
		}

		for _, v := range portfolioPositions {
			if v.InstrumentType == "bond" {
				operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.AssetUid, account.ID)
				if err != nil {
					return err
				}
				resultBondPosition, err := c.ProcessOperations(ctx, operationsDb)
				if err != nil {
					return err
				}
				totalAmount := portfolio.TotalAmount.ToFloat()

				bondReport, err := c.CreateGeneralBondReport(ctx, resultBondPosition, totalAmount)
				if err != nil {
					return err
				}
				switch {
				case bondReport.Replaced:
					tickerTimeKey := domain.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.ReplacedBondsReport[tickerTimeKey] = bondReport
				case bondReport.Currencies != "rub":
					tickerTimeKey := domain.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.EuroBondsReport[tickerTimeKey] = bondReport
				default:
					tickerTimeKey := domain.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.RubBondsReport[tickerTimeKey] = bondReport
				}

			}
		}

		err = Vizualization(logg, &generalBondReports, chatID, account.ID)
		if err != nil {
			return err
		}

	}

	return nil
}

func Vizualization(logger *slog.Logger, generalBondReports *domain.GeneralBondReports, chatID int, accountID string) (err error) {
	const op = "service.Vizualization"

	start := time.Now()
	logg := logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't do vizualization", err)
	}()
	reports := make([][]domain.GeneralBondReportPosition, 0)

	rubbleBondReportSorted, err := sortGeneralBondReports(logg, generalBondReports.RubBondsReport)
	if err != nil && !errors.Is(err, domain.ErrEmptyReport) {
		return err
	}
	replacedBondReportSorted, err := sortGeneralBondReports(logg, generalBondReports.ReplacedBondsReport)
	if err != nil && !errors.Is(err, domain.ErrEmptyReport) {
		return err
	}
	euroBondReportSorted, err := sortGeneralBondReports(logg, generalBondReports.EuroBondsReport)
	if err != nil && !errors.Is(err, domain.ErrEmptyReport) {
		return err
	}
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
			typeOfBonds = domain.ReplacedBonds
		case report[0].Currencies != "rub":
			typeOfBonds = domain.EuroBonds
		default:
			typeOfBonds = domain.RubBonds
		}
		pathDir := path.Join(reportPath, strconv.Itoa(chatID), accountID)
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
			err := visualization.Vizualize(logg, report[start:end], pathAndName, typeOfBonds)
			if err != nil {
				return e.WrapIfErr("vizualize error", err)
			}
			count += 1

		}
	}
	return nil
}

func (c *Client) GetBondReports(ctx context.Context, chatID int) (_ domain.BondReportsResponce, err error) {
	const op = "service.GetBondReports"
	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get general bond report", err)
	}()

	reportsInByteByAccounts := make([][]*domain.MediaGroup, 0)

	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return domain.BondReportsResponce{}, err
	}

	for _, account := range accounts {
		err = c.updateOperations(ctx, chatID, account.ID, account.OpenedDate)
		if err != nil {
			return domain.BondReportsResponce{}, err
		}
		if account.Status != 2 {
			continue
		}
		portfolio, err := c.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return domain.BondReportsResponce{}, err
		}

		portfolioPositions, err := c.MapPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
		if err != nil {
			return domain.BondReportsResponce{}, err
		}

		err = c.Storage.DeleteGeneralBondReport(context.Background(), chatID, account.ID)
		if err != nil {
			return domain.BondReportsResponce{}, err
		}
		generalBondReports := domain.GeneralBondReports{
			RubBondsReport:      make(map[domain.TickerTimeKey]domain.GeneralBondReportPosition),
			EuroBondsReport:     make(map[domain.TickerTimeKey]domain.GeneralBondReportPosition),
			ReplacedBondsReport: make(map[domain.TickerTimeKey]domain.GeneralBondReportPosition),
		}

		for _, v := range portfolioPositions {
			if v.InstrumentType == "bond" {
				operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.AssetUid, account.ID)
				if err != nil {
					return domain.BondReportsResponce{}, err
				}
				resultBondPosition, err := c.ProcessOperations(ctx, operationsDb)
				if err != nil {
					return domain.BondReportsResponce{}, err
				}
				totalAmount := portfolio.TotalAmount.ToFloat()

				bondReport, err := c.CreateGeneralBondReport(ctx, resultBondPosition, totalAmount)
				if err != nil {
					return domain.BondReportsResponce{}, err
				}
				switch {
				case bondReport.Replaced:
					tickerTimeKey := domain.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.ReplacedBondsReport[tickerTimeKey] = bondReport
				case bondReport.Currencies != "rub":
					tickerTimeKey := domain.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.EuroBondsReport[tickerTimeKey] = bondReport
				default:
					tickerTimeKey := domain.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.RubBondsReport[tickerTimeKey] = bondReport
				}

			}
		}

		reportsInByte, err := c.PrepareToGenerateTablePNG(&generalBondReports, chatID, account.ID)
		if err != nil {
			return domain.BondReportsResponce{}, err
		}
		reportsInByteByAccounts = append(reportsInByteByAccounts, reportsInByte)

	}
	getBondReportsResponce := domain.BondReportsResponce{Media: reportsInByteByAccounts}
	return getBondReportsResponce, nil
}

func (c *Client) PrepareToGenerateTablePNG(generalBondReports *domain.GeneralBondReports, chatID int, accountID string) (_ []*domain.MediaGroup, err error) {
	const op = "service.PrepareToGenerateTablePNG"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't prepareToGeneratePNG", err)
	}()

	reports := make([][]domain.GeneralBondReportPosition, 0)

	rubbleBondReportSorted, err := sortGeneralBondReports(logg, generalBondReports.RubBondsReport)
	if err != nil && !errors.Is(err, domain.ErrEmptyReport) {
		return nil, err
	}
	replacedBondReportSorted, err := sortGeneralBondReports(logg, generalBondReports.ReplacedBondsReport)
	if err != nil && !errors.Is(err, domain.ErrEmptyReport) {
		return nil, err
	}
	euroBondReportSorted, err := sortGeneralBondReports(logg, generalBondReports.EuroBondsReport)
	if err != nil && !errors.Is(err, domain.ErrEmptyReport) {
		return nil, err
	}
	reports = append(reports, rubbleBondReportSorted)
	reports = append(reports, replacedBondReportSorted)
	reports = append(reports, euroBondReportSorted)
	reportsInByte := make([]*domain.MediaGroup, 3)
	for i, report := range reports {
		reportsInByte[i] = domain.NewMediaGroup()
		mediaGroup := reportsInByte[i]
		if len(report) == 0 {
			continue
		}

		var typeOfBonds string
		switch {
		case report[0].Replaced:
			typeOfBonds = domain.ReplacedBonds
		case report[0].Currencies != "rub":
			typeOfBonds = domain.EuroBonds
		default:
			typeOfBonds = domain.RubBonds
		}
		count := 1
		for start := 0; start < len(report); start += 10 {
			end := start + 10
			if end > len(report) {
				end = len(report)
			}
			pngData, err := visualization.GenerateTablePNG(logg, report[start:end], typeOfBonds)
			if err != nil {
				return nil, e.WrapIfErr("vizualize error", err)
			}
			imageData := domain.NewImageData()
			imageData.Name = fmt.Sprintf("file%s_%v", typeOfBonds, count)
			imageData.Data = pngData
			imageData.Caption = typeOfBonds

			mediaGroup.Reports = append(mediaGroup.Reports, imageData)
			count += 1
		}
	}
	return reportsInByte, nil
}

func sortGeneralBondReports(logger *slog.Logger, report map[domain.TickerTimeKey]domain.GeneralBondReportPosition) (_ []domain.GeneralBondReportPosition, err error) {
	const op = "service.sortGeneralBondReports"

	start := time.Now()
	logg := logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't sort general bond report", err)
	}()

	// TODO: обработать более читаемо и в дальнейшем проверять ошибку
	if len(report) == 0 {
		return []domain.GeneralBondReportPosition{}, domain.ErrEmptyReport
	}

	keys := make([]domain.TickerTimeKey, 0, len(report))
	for k := range report {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Time.Equal(keys[j].Time) {
			return keys[i].Ticker < keys[j].Ticker
		}
		return keys[i].Time.Before(keys[j].Time)
	})
	result := make([]domain.GeneralBondReportPosition, len(keys))
	for i, k := range keys {
		result[i] = report[k]
	}

	return result, nil
}

func (c *Client) GetAccountsList(ctx context.Context) (answ domain.AccountListResponce, err error) {
	const op = "service.GetAccountsList"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get accounts", err)
	}()

	var accStr string = "По данному аккаунту доступны следующие счета:"

	accs, err := c.GetAccounts(ctx)
	if err != nil {
		return domain.AccountListResponce{}, err
	}
	for _, account := range accs {
		accStr += fmt.Sprintf("\n ID:%s, Type: %s, Name: %s, Status: %v \n", account.ID, account.Type, account.Name, account.Status)
	}
	accountResponce := domain.AccountListResponce{Accounts: accStr}
	return accountResponce, nil
}

func (c *Client) GetUsd(ctx context.Context) (_ domain.UsdResponce, err error) {
	const op = "service.GetUsd"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("usd get error", err)
	}()

	usd, err := c.GetCurrencyFromCB(ctx, "usd", time.Now())
	if err != nil {
		return domain.UsdResponce{}, err
	}
	usdResponce := domain.UsdResponce{Usd: usd}

	return usdResponce, nil
}

func (c *Client) updateOperations(ctx context.Context, chatID int, accountID string, openDate time.Time) (err error) {
	const op = "service.updateOperations"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't updateOperations", err)
	}()

	fromDate, err := c.Storage.LastOperationTime(context.Background(), chatID, accountID)
	// TODO: Если fromDate будет больше time.Now, то будет ошибка.
	// Есть вероятность такой ошибки, поэтому при тестировании функции нужно придумать другой способ вызова функции по последней операции
	fromDate = fromDate.Add(time.Microsecond * 1)

	if err != nil {
		if errors.Is(err, domain.ErrNoOpperations) {
			fromDate = openDate
		} else {
			return err
		}
	}

	tinkoffOperations, err := c.TinkoffGetOperations(ctx, accountID, fromDate)
	if err != nil {
		return err
	}

	operations := mapper.MapOperationToOperationWithoutCustomTypes(tinkoffOperations)

	err = c.Storage.SaveOperations(context.Background(), chatID, accountID, operations)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccounts(ctx context.Context) (_ map[string]domain.Account, err error) {
	const op = "service.GetAccounts"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("cant' get accounts", err)
	}()

	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (c *Client) GetPortfolioStructureForEachAccount(ctx context.Context) (_ domain.PortfolioStructureForEachAccountResponce, err error) {
	const op = "service.GetPortfolioStructureForEachAccount"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("cant' Get Portfolio Structure For Each Account", err)
	}()

	accounts, err := c.GetAccounts(ctx)
	response := domain.PortfolioStructureForEachAccountResponce{}
	if err != nil {
		return domain.PortfolioStructureForEachAccountResponce{}, err
	}
	for _, account := range accounts {
		if account.Status == 3 {
			continue
		}
		report, err := c.getPortfolioStructure(ctx, account)
		if err != nil {
			return domain.PortfolioStructureForEachAccountResponce{}, err
		}
		response.PortfolioStructures = append(response.PortfolioStructures, report)
	}
	return response, nil
}

func (c *Client) getPortfolioStructure(ctx context.Context, account domain.Account) (_ string, err error) {
	const op = "service.getPortfolioStructure"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("cant' get portfolio structure", err)
	}()

	portfolio, err := c.TinkoffGetPortfolio(ctx, account)
	if err != nil {
		return "", err
	}
	positions := portfolio.Positions

	accountTitle := fmt.Sprintf("Струтура брокерского счета: %s\n", account.Name)
	potfolioStructure, err := c.DivideByType(ctx, positions)
	if err != nil {
		return "", err
	}
	respPotfolioStructure := c.ResponsePortfolioStructure(potfolioStructure)

	response := accountTitle + respPotfolioStructure
	return response, nil
}

func (c *Client) GetUnionPortfolioStructureForEachAccount(ctx context.Context) (_ domain.UnionPortfolioStructureResponce, err error) {
	const op = "service.GetUnionPortfolioStructureForEachAccount"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("cant' Get Union Portfolio Structure For Each Account", err)
	}()

	accounts, err := c.GetAccounts(ctx)
	response := domain.UnionPortfolioStructureResponce{}
	if err != nil {
		return domain.UnionPortfolioStructureResponce{}, err
	}
	unionPortfolioStructure, err := c.getUnionPortfolioStructure(ctx, accounts)
	if err != nil {
		return domain.UnionPortfolioStructureResponce{}, err
	}
	response.Report = unionPortfolioStructure

	return response, nil
}

func (c *Client) getUnionPortfolioStructure(ctx context.Context, accounts map[string]domain.Account) (_ string, err error) {
	const op = "service.getUnionPortfolioStructure"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get union portfolio structuret", err)
	}()

	positionsList := make([]*domain.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		if account.Status != 2 {
			continue
		}
		portfolio, err := c.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return "", err
		}
		positions := portfolio.Positions

		potfolioStructure, err := c.DivideByType(ctx, positions)
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

func (c *Client) GetUnionPortfolioStructureWithSber(ctx context.Context) (_ domain.UnionPortfolioStructureWithSberResponce, err error) {
	const op = "service.GetUnionPortfolioStructureWithSber"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get union portfolio structure with Sber", err)
	}()

	responce := domain.UnionPortfolioStructureWithSberResponce{}
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, err
	}
	positionsList := make([]*domain.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		if account.Status != 2 {
			continue
		}
		portfolio, err := c.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return domain.UnionPortfolioStructureWithSberResponce{}, err
		}
		positions := portfolio.Positions

		potfolioStructure, err := c.DivideByType(ctx, positions)
		if err != nil {
			return domain.UnionPortfolioStructureWithSberResponce{}, err
		}
		positionsList = append(positionsList, potfolioStructure)
	}

	sberPortfolio, err := c.DivideByTypeFromSber(ctx, c.External.Sber.Portfolio)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, err
	}

	positionsList = append(positionsList, sberPortfolio)

	accountTitle := "Струтура всех инвестиций\n"
	unionPositions, err := c.UnionPortf(positionsList)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, err
	}
	vizualizeUnionPositions := c.ResponsePortfolioStructure(unionPositions)
	out := accountTitle + vizualizeUnionPositions
	responce.Report = out
	return responce, nil
}

func (c *Client) DivideByType(ctx context.Context, positions []domain.PortfolioPosition) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.DivideByType"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't divide by type", err)
	}()

	portfolio := domain.NewPortfolioByTypeAndCurrency()
	date := time.Now()

	if len(positions) == 0 {
		return portfolio, errors.New("positions are empty")
	}

	for _, pos := range positions {
		var positionPrice float64
		currencyOfPos := pos.CurrentPrice.Currency
		vunit_rate := 1.0
		if currencyOfPos != futuresPt && currencyOfPos != rub {
			vunit_rate, err = c.GetCurrencyFromCB(ctx, currencyOfPos, date)
			if err != nil {
				return portfolio, e.WrapIfErr("can't divide by type", err)
			}
		}
		positionPrice = pos.Quantity.ToFloat() * pos.CurrentPrice.ToFloat() * vunit_rate

		switch pos.InstrumentType {
		case bond:
			positionPrice += pos.CurrentNkd.ToFloat() * pos.Quantity.ToFloat() * vunit_rate
			portfolio.BondsAssets.SumOfAssets += positionPrice
			if _, exist := portfolio.BondsAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.BondsAssets.AssetsByCurrency[currencyOfPos] = domain.NewAssetsByParam()
			}
			portfolio.BondsAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case share:
			portfolio.SharesAssets.SumOfAssets += positionPrice
			if _, exist := portfolio.SharesAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.SharesAssets.AssetsByCurrency[currencyOfPos] = domain.NewAssetsByParam()
			}
			portfolio.SharesAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case futures:
			futures, err := c.TinkoffGetFutureBy(ctx, pos.Figi)
			if err != nil {
				return portfolio, e.WrapIfErr("can't get future data", err)
			}

			positionPrice = positionPrice / futures.MinPriceIncrement.ToFloat() * futures.MinPriceIncrementAmount.ToFloat()
			portfolio.FuturesAssets.SumOfAssets += positionPrice

			futureType := futures.AssetType
			switch futureType {
			case commodityType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name] = domain.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Commodity.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice
			case currencyType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name] = domain.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Currency.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice
			case securityType:
				resp, err := c.TinkoffGetBaseShareFutureValute(ctx, futures.BasicAssetPositionUid)
				if err != nil {
					return nil, e.WrapIfErr("can't divide by type", err)
				}
				valute := resp.Currency
				if _, exist := portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute]; !exist {
					portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute] = domain.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Security.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute].SumOfAssets += positionPrice
			case indexType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name] = domain.NewAssetsByParam()
				}

				portfolio.FuturesAssets.AssetsByType.Index.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice

			}
			// Чтобы сумма фьюча не сумировалась с суммой всех активов, так как фактически я за тело фьючерса не заплатил
			positionPrice = 0

		case etf:
			portfolio.EtfsAssets.SumOfAssets += positionPrice

			if _, exist := portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos] = domain.NewAssetsByParam()
			}
			portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case currency:
			curr, err := c.TinkoffGetCurrencyBy(ctx, pos.Figi)
			if err != nil {
				return portfolio, e.WrapIfErr("can't divide by type", err)
			}
			currName := curr.Isin
			portfolio.CurrenciesAssets.SumOfAssets += positionPrice

			if _, exist := portfolio.CurrenciesAssets.AssetsByCurrency[currName]; !exist {
				portfolio.CurrenciesAssets.AssetsByCurrency[currName] = domain.NewAssetsByParam()
			}
			portfolio.CurrenciesAssets.AssetsByCurrency[currName].SumOfAssets += positionPrice

		default:
		}
		portfolio.AllAssets += positionPrice
	}

	return portfolio, nil
}

func (c *Client) DivideByTypeFromSber(ctx context.Context, positions map[string]float64) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.DivideByTypeFromSber"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't divide by type from sber", err)
	}()

	portfolio := domain.NewPortfolioByTypeAndCurrency()

	if len(positions) == 0 {
		return portfolio, errors.New("positions are empty")
	}
	for ticker, quantity := range positions {
		positionsClassCodeVariants, err := c.TinkoffFindBy(ctx, ticker)
		if err != nil {
			return nil, e.WrapIfErr("can't divide by type from sber", err)
		}
		if len(positionsClassCodeVariants) == 0 {
			return nil, errors.New("positions variants are empty")
		}

		switch positionsClassCodeVariants[0].InstrumentType {
		case bond:
			bondUid := positionsClassCodeVariants[0].Uid
			bond, err := c.TinkoffGetBondByUid(ctx, bondUid)
			if err != nil {
				return nil, e.WrapIfErr("can't divide by type from sber", err)
			}
			currentNkd := bond.AciValue.ToFloat()
			currency := bond.Currency
			resp, err := c.TinkoffGetLastPriceInPersentageToNominal(ctx, bondUid)
			if err != nil {
				return nil, e.WrapIfErr("can't divide by type from sber", err)
			}
			currentPriceInPers := resp.LastPrice.ToFloat()
			currentPrice := currentPriceInPers / 100 * bond.Nominal.ToFloat()
			currentNkdOfPosition := currentNkd * quantity
			positionPrice := currentPrice*quantity + currentNkdOfPosition

			portfolio.AllAssets += positionPrice
			portfolio.BondsAssets.SumOfAssets += positionPrice

			if existing, exist := portfolio.BondsAssets.AssetsByCurrency[currency]; !exist {
				portfolio.BondsAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
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

func (c *Client) ResponsePortfolioStructure(portfolio *domain.PortfolioByTypeAndCurrency) string {
	const op = "service.ResponsePortfolioStructure"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
		)
	}()

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

func (c *Client) UnionPortf(portfolios []*domain.PortfolioByTypeAndCurrency) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.UnionPortf"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Info(fmt.Sprintf("start %s", op))
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't union portfolios ", err)
	}()

	unionPortf := domain.NewPortfolioByTypeAndCurrency()
	for _, portf := range portfolios {
		unionPortf.AllAssets += portf.AllAssets

		unionPortf.BondsAssets.SumOfAssets += portf.BondsAssets.SumOfAssets
		for k, v := range portf.BondsAssets.AssetsByCurrency {
			if existing, exist := unionPortf.BondsAssets.AssetsByCurrency[k]; !exist {
				unionPortf.BondsAssets.AssetsByCurrency[k] = domain.NewAssetsByParam()
				unionPortf.BondsAssets.AssetsByCurrency[k] = v
			} else {
				existing.SumOfAssets += v.SumOfAssets
			}
		}

		unionPortf.SharesAssets.SumOfAssets += portf.SharesAssets.SumOfAssets
		for currency, asset := range portf.SharesAssets.AssetsByCurrency {
			if existing, exists := unionPortf.SharesAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.SharesAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.EtfsAssets.SumOfAssets += portf.EtfsAssets.SumOfAssets
		for currency, asset := range portf.EtfsAssets.AssetsByCurrency {
			if existing, exists := unionPortf.EtfsAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.EtfsAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.CurrenciesAssets.SumOfAssets += portf.CurrenciesAssets.SumOfAssets
		for currency, asset := range portf.CurrenciesAssets.AssetsByCurrency {
			if existing, exists := unionPortf.CurrenciesAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.CurrenciesAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
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
				unionPortf.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Currency.SumOfAssets += portf.FuturesAssets.AssetsByType.Currency.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Security.SumOfAssets += portf.FuturesAssets.AssetsByType.Security.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Security.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Security.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Security.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Index.SumOfAssets += portf.FuturesAssets.AssetsByType.Index.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Index.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Index.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Index.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}
	}
	return unionPortf, nil
}
