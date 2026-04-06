package app

import (
	"bonds-report-service/internal/adapters/outbound/tinkoffApi/client/analyticsclient"
	"bonds-report-service/internal/adapters/outbound/tinkoffApi/client/instrumentsclient"
	"bonds-report-service/internal/adapters/outbound/tinkoffApi/client/portfolioclient"
	tinkofftransport "bonds-report-service/internal/adapters/outbound/tinkoffApi/transport"
	bondreport "bonds-report-service/internal/application/helpers/bondReport"
	cbrHelper "bonds-report-service/internal/application/helpers/cbr"
	dividerbyassettype "bonds-report-service/internal/application/helpers/dividerByAssetType"
	generalbondreport "bonds-report-service/internal/application/helpers/generalBondReport"
	moexHelper "bonds-report-service/internal/application/helpers/moex"
	positionProcessor "bonds-report-service/internal/application/helpers/positionsToPositionsWithAssetUidProcessor"
	"bonds-report-service/internal/application/helpers/report"
	reportlinebuiler "bonds-report-service/internal/application/helpers/reportLineBuiler"
	tinkoffHelper "bonds-report-service/internal/application/helpers/tinkoff"
	"bonds-report-service/internal/application/helpers/uidprovider"
	updateoperations "bonds-report-service/internal/application/helpers/updateOperations"
	"bonds-report-service/internal/application/ports"
	"log/slog"
)

func InitBondReportProcessor(logger *slog.Logger) *bondreport.BondReporter {
	logger.Info("initialize bond report processor")
	bondReporter := bondreport.NewBondReporter(logger)
	return bondReporter
}

func InitCBRCurrencyGetter(logger *slog.Logger, cbr ports.CbrClient, storage ports.Storage) *cbrHelper.CbrHelper {
	logger.Info("initialize cbr currency getter")
	cbrCurrencyGetter := cbrHelper.NewCbrHelper(logger, cbr, storage)
	return cbrCurrencyGetter
}

func InitGeneralReportProcessor(logger *slog.Logger) *generalbondreport.GeneralBondReporter {
	logger.Info("initialize general bond report processor")
	generalBondReporter := generalbondreport.NewGeneralBondReporter(logger)
	return generalBondReporter
}

func InitMoexSpecificationGetter(logger *slog.Logger, moex ports.MoexClient) *moexHelper.MoexHelper {
	logger.Info("initialize moex specification getter")
	moexSpecificationGetter := moexHelper.NewMoexHelper(logger, moex)
	return moexSpecificationGetter
}

func InitReportProcessor(logger *slog.Logger) *report.ReportProcessor {
	logger.Info("initialize report processor")
	reportProcessor := report.NewReportProcessor(logger)
	return reportProcessor
}

func InitTinkoffApiHelper(logger *slog.Logger, host string) *tinkoffHelper.TinkoffHelper {
	logger.Info("initialize Tinkoff helper", slog.String("address", host))
	if host == "" {
		panic("tinkoff host is empty")
	}
	transport := tinkofftransport.NewTransport(logger, host)
	analyticsclient := analyticsclient.NewAnalyticsTinkoffClient(logger, transport)
	instrumentsclient := instrumentsclient.NewInstrumentsTinkoffClient(logger, transport)
	portfolioclient := portfolioclient.NewPortfolioTinkoffClient(logger, transport)
	helper := tinkoffHelper.NewTinkoffHelper(logger, instrumentsclient, portfolioclient, analyticsclient)
	return helper
}

func InitUidProvider(logger *slog.Logger, repo ports.Storage, analyticService ports.TinkoffAnalyticsClient) *uidprovider.UidProvider {
	logger.Info("initialize uid provider")
	uidProvider := uidprovider.NewUidProvider(repo, analyticService)
	return uidProvider
}

func InitOperationsUpdater(logger *slog.Logger, tinkoffHelper *tinkoffHelper.TinkoffHelper, storage ports.Storage) *updateoperations.Updater {
	logger.Info("initialize operations updater")
	operrationsUpdater := updateoperations.NewUpdater(logger, storage, tinkoffHelper)
	return operrationsUpdater
}

func InitPositionProcessor(logger *slog.Logger, uidProvider ports.UidProvider) *positionProcessor.Processor {
	logger.Info("initialize position processor")
	positionProcessor := positionProcessor.NewProcessor(logger, uidProvider)
	return positionProcessor
}

func InitReportLineBuilder(logger *slog.Logger, tinkoffHelper *tinkoffHelper.TinkoffHelper, cbrHelper ports.CbrCurrencyGetter) *reportlinebuiler.ReportLineBuilder {
	logger.Info("initialize report line builder")
	reportlinebuiler := reportlinebuiler.NewReportLineBuilder(logger, tinkoffHelper, cbrHelper)
	return reportlinebuiler
}

func InitDividerByAssetType(logger *slog.Logger, tinkoffHelper *tinkoffHelper.TinkoffHelper, cbrHelper ports.CbrCurrencyGetter, workersNumber int) *dividerbyassettype.DividerByAssetType {
	logger.Info("initialize divider by asset type")
	dividerByAssetType := dividerbyassettype.NewDividerByAssetType(logger, tinkoffHelper, cbrHelper, workersNumber)
	return dividerByAssetType
}
