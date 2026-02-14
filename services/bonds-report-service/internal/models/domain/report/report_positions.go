package report

import (
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"log/slog"
	"time"
)

const (
	EuroTransBuyCost = 240   // Стоимость Евротранса при переводе из другого депозитария
	threeYearInHours = 26304 // Три года в часах
	baseTaxRate      = 0.13  // Налог с продажи ЦБ
)

type ReportPositions struct {
	Quantity         float64
	CurrentPositions []PositionByFIFO
}

func NewReportPositons() *ReportPositions {
	return &ReportPositions{
		Quantity:         0,
		CurrentPositions: []PositionByFIFO{},
	}
}

type PositionByFIFO struct {
	Name                  string
	Replaced              bool
	CurrencyIfReplaced    string
	BuyDate               time.Time
	SellDate              time.Time
	Quantity              float64
	Figi                  string
	InstrumentType        string
	InstrumentUid         string
	Ticker                string
	ClassCode             string
	Nominal               float64
	BuyPrice              float64
	SellPrice             float64 // Для открытых позиций.Текущая цена с биржи
	BuyPayment            float64
	SellPayment           float64
	Currency              string
	BuyAccruedInt         float64 // НКД при покупке
	SellAccruedInt        float64
	PartialEarlyRepayment float64 // Частичное досрочное гашение
	TotalCoupon           float64
	TotalDividend         float64
	TotalComission        float64
	PaidTax               float64 // Фактически уплаченный налог(Часть налога будет уплачена в конце года, либо при выводе средств)
	TotalTax              float64 // Налог рассчитанный
	PositionProfit        float64 // С учетом рассчитанных налогов(TotalTax)
	ProfitInPercentage    float64 // В процентах строковая переменная
}

func (p *PositionByFIFO) ApplyBondMetadata(
	bond domain.BondIdentIdentifiers,
) {
	p.Ticker = bond.Ticker
	p.ClassCode = bond.ClassCode
	p.Replaced = bond.Replaced
	p.CurrencyIfReplaced = bond.NominalCurrency
}

// 2	Удержание НДФЛ по купонам.
// 8    Удержание налога по дивидендам.
func (p *ReportPositions) ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(
	ctx context.Context,
	logger *slog.Logger,
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "service.processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends"

	defer logging.LogOperation_Debug(ctx, logger, op, &err)()

	if p.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range p.CurrentPositions {
			currentPosition := &p.CurrentPositions[i]
			proportion := currentPosition.Quantity / p.Quantity
			currentPosition.PaidTax += operation.Payment * proportion
		}
	}
	return nil
}

// 10	Частичное погашение облигаций.
func (p *ReportPositions) ProcessPartialRedemptionOfBonds(
	ctx context.Context,
	logger *slog.Logger,
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "service.processPartialRedemptionOfBonds"

	defer logging.LogOperation_Debug(ctx, logger, op, &err)()

	if p.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range p.CurrentPositions {
			currentPosition := &p.CurrentPositions[i]
			proportion := currentPosition.Quantity / p.Quantity
			currentPosition.PartialEarlyRepayment += operation.Payment * proportion
		}
	}
	return nil
}

// 15	Покупка ЦБ.
// 16	Покупка ЦБ с карты.
// 57   Перевод ценных бумаг с ИИС на Брокерский счет
func (p *ReportPositions) ProcessPurchaseOfSecurities(ctx context.Context,
	logger *slog.Logger,
	operation domain.OperationWithoutCustomTypes,
	bondIdentifiers domain.BondIdentIdentifiers,
	lastPrice domain.LastPrice,
	vunitRate domain.Rate,
) {
	const op = "service.processPurchaseOfSecurities"

	defer logging.LogOperation_Debug(ctx, logger, op, nil)()

	// при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := PositionByFIFO{
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

	position.ApplyBondMetadata(bondIdentifiers)

	nominal := CalculateNominal(bondIdentifiers.Nominal, position.Replaced, vunitRate)

	sellPrice := CalculateSellPrice(nominal, lastPrice)

	position.Nominal = nominal
	position.SellPrice = sellPrice

	p.CurrentPositions = append(p.CurrentPositions, position)
	p.Quantity += operation.QuantityDone
}

// 17	Перевод ценных бумаг из другого депозитария.
func (p *ReportPositions) ProcessTransferOfSecuritiesFromAnotherDepository(
	ctx context.Context,
	logger *slog.Logger,
	operation domain.OperationWithoutCustomTypes,
	bondIdentifiers domain.BondIdentIdentifiers,
	lastPrice domain.LastPrice,
	vunitRate domain.Rate,
) {
	const op = "service.processTransferOfSecuritiesFromAnotherDepository"

	defer logging.LogOperation_Debug(ctx, logger, op, nil)()
	// при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := PositionByFIFO{
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

	position.ApplyBondMetadata(bondIdentifiers)

	nominal := CalculateNominal(bondIdentifiers.Nominal, position.Replaced, vunitRate)

	sellPrice := CalculateSellPrice(nominal, lastPrice)

	position.Nominal = nominal
	position.SellPrice = sellPrice

	p.CurrentPositions = append(p.CurrentPositions, position)
	p.Quantity += operation.QuantityDone
}

// 21	Выплата дивидендов.
func (p *ReportPositions) ProcessPaymentOfDividends(
	ctx context.Context,
	logger *slog.Logger,
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "service.processPaymentOfDividends"

	defer logging.LogOperation_Debug(ctx, logger, op, &err)()

	if p.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range p.CurrentPositions {
			currentPosition := &p.CurrentPositions[i]
			proportion := currentPosition.Quantity / p.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalDividend += operation.Payment * proportion
		}
	}
	return nil
}

// 23 Выплата купонов.
func (p *ReportPositions) ProcessPaymentOfCoupons(
	ctx context.Context,
	logger *slog.Logger,
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "service.processPaymentOfCoupons"

	defer logging.LogOperation_Debug(ctx, logger, op, &err)()

	if p.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range p.CurrentPositions {
			currentPosition := &p.CurrentPositions[i]
			proportion := currentPosition.Quantity / p.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalCoupon += operation.Payment * proportion
		}
	}
	return nil
}

// 47	Гербовый сбор.
func (p *ReportPositions) ProcessStampDuty(
	ctx context.Context,
	logger *slog.Logger,
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "service.processStampDuty"

	defer logging.LogOperation_Debug(ctx, logger, op, &err)()

	if p.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range p.CurrentPositions {
			currentPosition := &p.CurrentPositions[i]
			proportion := currentPosition.Quantity / p.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalComission += operation.Payment * proportion
		}
	}
	return nil
}

// 22	Продажа ЦБ.
func (p *ReportPositions) ProcessSellOfSecurities(
	ctx context.Context,
	logger *slog.Logger,
	operation *domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "service.processSellOfSecurities"

	defer logging.LogOperation_Debug(ctx, logger, op, &err)()

	p.Quantity -= operation.QuantityDone
	// TODO: Переписать ПОЗЖЕ Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления в которых Кол-во проданных
	// бумаг больше кол-ва бумаг в текущей позиции
	var deleteCount int
end:
	for i := range p.CurrentPositions {
		currPosition := &p.CurrentPositions[i]
		currentQuantity := currPosition.Quantity
		sellQuantity := operation.QuantityDone
		switch {
		case currentQuantity > sellQuantity:
			err := isCurrentQuantityGreaterThanSellQuantity(ctx, operation, currPosition)
			if err != nil {
				return err
			}
			// Прерываем цикл
			break end
		case currPosition.Quantity == operation.QuantityDone:
			err := isEqualCurrentQuantityAndSellQuantity(p)
			if err != nil {
				return err
			}
			// Прерываем цикл
			break end
		case currentQuantity < sellQuantity:
			// Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления
			deleteCount += 1
			err := isCurrentQuantityLessThanSellQuantity(ctx, operation, currPosition)
			if err != nil {
				return err
			}
		}

	}
	// удаляем закрытые позиции из среза текущих позиций
	p.CurrentPositions = p.CurrentPositions[deleteCount:]
	return nil
}

func isCurrentQuantityGreaterThanSellQuantity(
	ctx context.Context,
	operation *domain.OperationWithoutCustomTypes,
	currPosition *PositionByFIFO,
) error {
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

func isEqualCurrentQuantityAndSellQuantity(processPosition *ReportPositions) error {
	processPosition.CurrentPositions = processPosition.CurrentPositions[1:]
	return nil
}

func isCurrentQuantityLessThanSellQuantity(ctx context.Context, operation *domain.OperationWithoutCustomTypes, currPosition *PositionByFIFO) error {
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
