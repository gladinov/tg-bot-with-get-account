package dividerbyassettype

import (
	tinkoffHelper "bonds-report-service/internal/application/helpers/tinkoff"
	"bonds-report-service/internal/application/ports"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gladinov/e"
)

var ErrEmptyPosition = errors.New("positions are empty")

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

type DividerByAssetType struct {
	logger        *slog.Logger
	CbrGetter     ports.CbrCurrencyGetter
	TinkoffHelper *tinkoffHelper.TinkoffHelper
	now           func() time.Time
}

func NewDividerByAssetType(logger *slog.Logger, tinkoffHelper *tinkoffHelper.TinkoffHelper, cbrGetter ports.CbrCurrencyGetter) *DividerByAssetType {
	return &DividerByAssetType{
		logger:        logger,
		TinkoffHelper: tinkoffHelper,
		CbrGetter:     cbrGetter,
		now:           time.Now,
	}
}

func (d *DividerByAssetType) DivideByType(ctx context.Context, positions []domain.PortfolioPosition) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.DivideByType"

	defer logging.LogOperation_Debug(ctx, d.logger, op, &err)()

	portfolio := domain.NewPortfolioByTypeAndCurrency()
	date := d.now()

	if len(positions) == 0 {
		return portfolio, ErrEmptyPosition
	}

	for _, pos := range positions {
		var positionPrice float64
		currencyOfPos := pos.CurrentPrice.Currency

		var vunit_rate float64
		if currencyOfPos != futuresPt && currencyOfPos != rub {
			vunit_rate, err = d.CbrGetter.GetCurrencyFromCB(ctx, currencyOfPos, date)
			if err != nil {
				return nil, e.WrapIfErr("can't get currency from CB", err)
			}
		} else {
			vunit_rate = 1.0
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
			futures, err := d.TinkoffHelper.TinkoffGetFutureBy(ctx, pos.Figi)
			if err != nil {
				return nil, e.WrapIfErr("can't get future data", err)
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
				resp, err := d.TinkoffHelper.TinkoffGetBaseShareFutureValute(ctx, futures.BasicAssetPositionUid)
				if err != nil {
					return nil, e.WrapIfErr("can't dget base share future valute from tinkoff", err)
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
			curr, err := d.TinkoffHelper.TinkoffGetCurrencyBy(ctx, pos.Figi)
			if err != nil {
				return nil, e.WrapIfErr("can't get currency by figi from tinkoff", err)
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
