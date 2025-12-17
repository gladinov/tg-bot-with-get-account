package service

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"time"

	"bonds-report-service/internal/service/service_models"
	"bonds-report-service/lib/e"
)

const (
	WithholdingOfPersonalIncomeTaxOnCoupons        = 2     // 2	Удержание НДФЛ по купонам.
	WithholdingOfPersonalIncomeTaxOnDividends      = 8     // 8    Удержание налога по дивидендам.
	PartialRedemptionOfBonds                       = 10    // 10	Частичное погашение облигаций.
	PurchaseOfSecurities                           = 15    // 15	Покупка ЦБ.
	PurchaseOfSecuritiesWithACard                  = 16    // 16	Покупка ЦБ с карты.
	TransferOfSecuritiesFromAnotherDepository      = 17    // 17	Перевод ценных бумаг из другого депозитария.
	WithhouldingACommissionForTheTransaction       = 19    // 19	Удержание комиссии за операцию.
	PaymentOfDividends                             = 21    // 21	Выплата дивидендов.
	SaleOfSecurities                               = 22    // 22	Продажа ЦБ.
	PaymentOfCoupons                               = 23    // 23 Выплата купонов.
	StampDuty                                      = 47    // 47	Гербовый сбор.
	TransferOfSecuritiesFromIISToABrokerageAccount = 57    // 57   Перевод ценных бумаг с ИИС на Брокерский счет
	EuroTransBuyCost                               = 240   //Стоимость Евротранса при переводе из другого депозитария
	threeYearInHours                               = 26304 // Три года в часах
	baseTaxRate                                    = 0.13  // Налог с продажи ЦБ
)

func (c *Client) GetSpecificationsFromTinkoff(ctx context.Context, position *service_models.PositionByFIFO) (err error) {
	const op = "service.GetSpecificationsFromTinkoff"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not get specifications from tinkoff", err)
	}()

	resSpecFromTinkoff, err := c.TinkoffGetBondActions(ctx, position.InstrumentUid)
	if err != nil {
		return errors.New("service:GetSpecificationsFromMoex " + err.Error())
	}
	position.Ticker = resSpecFromTinkoff.Ticker
	position.ClassCode = resSpecFromTinkoff.ClassCode
	if resSpecFromTinkoff.Replaced {
		date := time.Now()
		position.Replaced = true
		isoCurrName := resSpecFromTinkoff.NominalCurrency
		position.CurrencyIfReplaced = isoCurrName
		vunit_rate, err := c.GetCurrencyFromCB(isoCurrName, date)
		if err != nil {
			return e.WrapIfErr("getSpecificationsFromMoex err", err)
		}
		position.Nominal = resSpecFromTinkoff.Nominal.ToFloat() * vunit_rate
	} else {
		position.Nominal = resSpecFromTinkoff.Nominal.ToFloat()
	}
	resp, err := c.TinkoffGetLastPriceInPersentageToNominal(ctx, position.InstrumentUid)
	if err != nil {
		return errors.New("service:GetSpecificationsFromMoex:" + err.Error())
	}
	resLastPriceFromTinkoff := resp.LastPrice.ToFloat()
	position.SellPrice = math.Round(resLastPriceFromTinkoff/100*position.Nominal*100) / 100
	return nil

}

func (c *Client) ProcessOperations(ctx context.Context, operations []service_models.Operation) (_ *service_models.ReportPositions, err error) {
	const op = "service.ProcessOperations"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process operations", err)
	}()

	processPosition := &service_models.ReportPositions{}
	for _, operation := range operations {
		switch operation.Type {
		// 2	Удержание НДФЛ по купонам.
		// 8    Удержание налога по дивидендам.
		case WithholdingOfPersonalIncomeTaxOnCoupons, WithholdingOfPersonalIncomeTaxOnDividends:
			err := c.processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(operation, processPosition)
			if err != nil {
				return nil, err
			}

			// 10	Частичное погашение облигаций.
		case PartialRedemptionOfBonds:
			err := c.processPartialRedemptionOfBonds(operation, processPosition)
			if err != nil {
				return nil, err
			}

			// 15	Покупка ЦБ.
			// 16	Покупка ЦБ с карты.
			// 57   Перевод ценных бумаг с ИИС на Брокерский счет
		case PurchaseOfSecurities, PurchaseOfSecuritiesWithACard, TransferOfSecuritiesFromIISToABrokerageAccount:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				err := c.processPurchaseOfSecurities(ctx, operation, processPosition)
				if err != nil {
					return nil, err
				}
			}
			// 17	Перевод ценных бумаг из другого депозитария.
		case TransferOfSecuritiesFromAnotherDepository:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				err := c.processTransferOfSecuritiesFromAnotherDepository(ctx, operation, processPosition)
				if err != nil {
					return nil, err
				}
			}
			// 19	Удержание комиссии за операцию.
		case WithhouldingACommissionForTheTransaction:
			// Посчитали комисссию в операции покупки(15,16.17,57) и операции продажи(22)

			// 21	Выплата дивидендов.
		case PaymentOfDividends:
			err := c.processPaymentOfDividends(operation, processPosition)
			if err != nil {
				return nil, err
			}
			// 22	Продажа ЦБ.
		case SaleOfSecurities:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				err := c.processSellOfSecurities(&operation, processPosition)
				if err != nil {
					return nil, err
				}
			}

			// 23 Выплата купонов.
		case PaymentOfCoupons:
			err := c.processPaymentOfCoupons(operation, processPosition)
			if err != nil {
				return nil, err
			}

			// 47	Гербовый сбор.
		case StampDuty:
			err := c.processStampDuty(operation, processPosition)
			if err != nil {
				return nil, err
			}
		default:
			continue

		}
	}
	return processPosition, nil
}

// 2	Удержание НДФЛ по купонам.
// 8    Удержание налога по дивидендам.
func (c *Client) processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(operation service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process withholding of personal income tax on coupons or dividends", err)
	}()
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			currentPosition.PaidTax += operation.Payment * proportion
		}
	}
	return nil
}

// 10	Частичное погашение облигаций.
func (c *Client) processPartialRedemptionOfBonds(operation service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processPartialRedemptionOfBonds"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process Partial Redemption Of Bonds", err)
	}()
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			currentPosition.PartialEarlyRepayment += operation.Payment * proportion
		}
	}
	return nil

}

// 15	Покупка ЦБ.
// 16	Покупка ЦБ с карты.
// 57   Перевод ценных бумаг с ИИС на Брокерский счет
func (c *Client) processPurchaseOfSecurities(ctx context.Context, operation service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processPurchaseOfSecurities"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process Purchase Of Securities", err)
	}()

	// при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := service_models.PositionByFIFO{
		Name:           operation.Name,
		BuyDate:        operation.Date,
		Figi:           operation.Figi,
		Quantity:       operation.QuantityDone,
		InstrumentType: operation.InstrumentType,
		InstrumentUid:  operation.InstrumentUid,
		BuyPrice:       operation.Price,
		Currency:       operation.Currency,
		BuyAccruedInt:  operation.AccruedInt, // НКД при покупке
		TotalComission: operation.Commission,
	}

	err = c.GetSpecificationsFromTinkoff(ctx, &position)
	if err != nil {
		return errors.New("service:processPurchaseOfSecurities:" + err.Error())
	}

	processPosition.CurrentPositions = append(processPosition.CurrentPositions, position)
	processPosition.Quantity += operation.QuantityDone
	return nil
}

// 17	Перевод ценных бумаг из другого депозитария.
func (c *Client) processTransferOfSecuritiesFromAnotherDepository(ctx context.Context, operation service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processTransferOfSecuritiesFromAnotherDepository"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process Transfer Of Securities From Another Depository", err)
	}()
	// при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := service_models.PositionByFIFO{
		Name:           operation.Name,
		BuyDate:        operation.Date,
		Figi:           operation.Figi,
		Quantity:       operation.QuantityDone,
		InstrumentType: operation.InstrumentType,
		InstrumentUid:  operation.InstrumentUid,
		BuyPrice:       operation.Price,
		Currency:       operation.Currency,
		BuyAccruedInt:  operation.AccruedInt, // НКД при покупке
		TotalComission: operation.Commission,
	}
	// Для Евротранса исключение
	if operation.InstrumentUid == "02b2ea14-3c4b-47e8-9548-45a8dbcc8f8a" {
		position.BuyPrice = EuroTransBuyCost
	}

	err = c.GetSpecificationsFromTinkoff(ctx, &position)
	if err != nil {
		return errors.New("service:processTransferOfSecuritiesFromAnotherDepository:" + err.Error())
	}

	processPosition.CurrentPositions = append(processPosition.CurrentPositions, position)
	processPosition.Quantity += operation.QuantityDone

	return nil
}

// 21	Выплата дивидендов.
func (c *Client) processPaymentOfDividends(operation service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processPaymentOfDividends"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process Payment Of Dividends", err)
	}()

	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalDividend += operation.Payment * proportion
		}
	}
	return nil
}

// 23 Выплата купонов.
func (c *Client) processPaymentOfCoupons(operation service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processPaymentOfCoupons"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process Payment Of Coupons", err)
	}()

	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalCoupon += operation.Payment * proportion
		}
	}
	return nil
}

// 47	Гербовый сбор.
func (c *Client) processStampDuty(operation service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processStampDuty"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process Payment Of Coupons", err)
	}()
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalComission += operation.Payment * proportion
		}
	}
	return nil
}

// 22	Продажа ЦБ.
func (c *Client) processSellOfSecurities(operation *service_models.Operation, processPosition *service_models.ReportPositions) (err error) {
	const op = "service.processSellOfSecurities"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("could not process Sell Of Securities", err)
	}()

	processPosition.Quantity -= operation.QuantityDone
	// Переписать ПОЗЖЕ Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления в которых Кол-во проданных
	// бумаг больше кол-ва бумаг в текущей позиции
	var deleteCount int
end:
	for i := range processPosition.CurrentPositions {
		currPosition := &processPosition.CurrentPositions[i]
		currentQuantity := currPosition.Quantity
		sellQuantity := operation.QuantityDone
		switch {
		case currentQuantity > sellQuantity:
			err := isCurrentQuantityGreaterThanSellQuantity(operation, currPosition)
			if err != nil {
				return err
			}
			// Прерываем цикл
			break end
		case currPosition.Quantity == operation.QuantityDone:
			err := isEqualCurrentQuantityAndSellQuantity(processPosition)
			if err != nil {
				return err
			}
			// Прерываем цикл
			break end
		case currentQuantity < sellQuantity:
			// Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления
			deleteCount += 1
			err := isCurrentQuantityLessThanSellQuantity(operation, currPosition)
			if err != nil {
				return err
			}
		}

	}
	// удаляем закрытые позиции из среза текущих позиций
	processPosition.CurrentPositions = processPosition.CurrentPositions[deleteCount:]
	return nil
}

func isCurrentQuantityGreaterThanSellQuantity(operation *service_models.Operation, currPosition *service_models.PositionByFIFO) error {
	currentQuantity := currPosition.Quantity
	sellQuantity := operation.QuantityDone
	var proportion float64
	if currentQuantity != 0 {
		proportion = sellQuantity / currentQuantity
	} else {
		return errors.New("divide by zero")
	}
	// Отнимаем кол-во проданных бумаг от количества бумаг в текущей позиции
	currPosition.Quantity -= sellQuantity
	// Изменяем значения текущей позиции, умножая на остаток от пропорции
	currPosition.TotalComission = currPosition.TotalComission * (1 - proportion)
	currPosition.PaidTax = currPosition.PaidTax * (1 - proportion)
	currPosition.BuyAccruedInt = currPosition.BuyAccruedInt * (1 - proportion)
	return nil
}

func isEqualCurrentQuantityAndSellQuantity(processPosition *service_models.ReportPositions) error {
	processPosition.CurrentPositions = processPosition.CurrentPositions[1:]
	return nil
}

func isCurrentQuantityLessThanSellQuantity(operation *service_models.Operation, currPosition *service_models.PositionByFIFO) error {
	currentQuantity := currPosition.Quantity
	sellQuantity := operation.QuantityDone
	var proportion float64
	if sellQuantity != 0 {
		proportion = currentQuantity / sellQuantity
	} else {
		return errors.New("divide by zero")
	}

	// НКД
	operation.AccruedInt -= operation.AccruedInt * proportion
	// Плюсуем комиссию за продажу бумаг
	operation.Commission -= operation.Commission * proportion

	// Изменяем значение Quantity.Operation
	operation.QuantityDone -= currPosition.Quantity

	return nil
}

// Доход по позиции до вычета налога
func getSecurityIncomeWithoutTax(p service_models.PositionByFIFO) float64 {
	// Для незакрытых позиций необходимо посчитать еще потенциальную комиссию при продаже
	quantity := p.Quantity
	buySellDifference := (p.SellPrice-p.BuyPrice)*quantity + p.SellAccruedInt - p.BuyAccruedInt
	cashFlow := p.TotalCoupon + p.TotalDividend
	positionProfit := buySellDifference + cashFlow + p.TotalComission + p.PartialEarlyRepayment
	return positionProfit
}

// Расход полного налога по закрытой позиции
func getTotalTaxFromPosition(p service_models.PositionByFIFO, profit float64) float64 {
	// Рассчитываем срок владения
	buyDate := p.BuyDate
	sellDate := p.SellDate
	timeDuration := sellDate.Sub(buyDate).Hours()
	var tax float64
	// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
	if profit > 0 && timeDuration < float64(threeYearInHours) {
		tax = profit * baseTaxRate
	} else {
		tax = 0
	}
	return tax
}

// Расчет прибыли после налогообложения
func getSecurityIncome(profit, tax float64) float64 {
	profitAfterTax := profit - tax
	return profitAfterTax
}

func getProfitInPercentage(profit, buyPrice, quantity float64) (float64, error) {
	if buyPrice != 0 || quantity != 0 {
		profitInPercentage := RoundFloat((profit/(buyPrice*quantity))*100, 2)
		return profitInPercentage, nil
	} else {
		return 0, errors.New("divide by zero")
	}
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
