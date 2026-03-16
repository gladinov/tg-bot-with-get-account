package report

import (
	"bonds-report-service/internal/domain"

	"github.com/gladinov/e"
)

// TODO: Подумать как разбить адекватно эту мешанину из файлов

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

func (p *ReportPositions) Apply(
	operation domain.OperationWithoutCustomTypes,
	bond domain.BondIdentIdentifiers,
	lastPrice domain.LastPrice,
	rate domain.Rate,
) error {
	switch operation.Type {
	// 2	Удержание НДФЛ по купонам.
	// 8    Удержание налога по дивидендам.
	case WithholdingOfPersonalIncomeTaxOnCoupons, WithholdingOfPersonalIncomeTaxOnDividends:
		if err := p.ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(
			operation); err != nil {
			return e.WrapIfErr("failed to Process Withholding Of Personal Income Tax On Coupons Or Dividends", err)
		}
		return nil

		// 10	Частичное погашение облигаций.
	case PartialRedemptionOfBonds:
		if err := p.ProcessPartialRedemptionOfBonds(
			operation); err != nil {
			return e.WrapIfErr("failed to Process Partial Redemption Of Bonds", err)
		}
		return nil

		// 15	Покупка ЦБ.
		// 16	Покупка ЦБ с карты.
		// 17	Перевод ценных бумаг из другого депозитария.
		// 57   Перевод ценных бумаг с ИИС на Брокерский счет
	case PurchaseOfSecurities,
		PurchaseOfSecuritiesWithACard,
		TransferOfSecuritiesFromIISToABrokerageAccount,
		TransferOfSecuritiesFromAnotherDepository:
		// Проверяем операцию на выполнение.
		// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
		if operation.QuantityDone == 0 {
			return ErrZeroQuantity
		}
		p.ProcessPurchaseOfSecurities(
			operation,
			bond,
			lastPrice,
			rate)
		return nil
		// 19	Удержание комиссии за операцию.
	case WithhouldingACommissionForTheTransaction:
		// Посчитали комисссию в операции покупки(15,16.17,57) и операции продажи(22)
		return nil
		// 21	Выплата дивидендов.
	case PaymentOfDividends:
		if err := p.ProcessPaymentOfDividends(operation); err != nil {
			return e.WrapIfErr("failed to Process Payment Of Dividends", err)
		}
		return nil
		// 22	Продажа ЦБ.
	case SaleOfSecurities:
		// Проверяем операцию на выполнение.
		// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
		if operation.QuantityDone == 0 {
			return ErrZeroQuantity
		} else {
			if err := p.ProcessSellOfSecurities(&operation); err != nil { // TODO: изменяем операцию . Надо здесь обдумать верна ли логика
				return e.WrapIfErr("failed to Process Sell Of Securities", err)
			}
			return nil
		}

		// 23 Выплата купонов.
	case PaymentOfCoupons:
		if err := p.ProcessPaymentOfCoupons(operation); err != nil {
			return e.WrapIfErr("failed to Process Payment Of Coupons", err)
		}
		return nil

		// 47	Гербовый сбор.
	case StampDuty:
		if err := p.ProcessStampDuty(operation); err != nil {
			return e.WrapIfErr("failed to Process Stamp Duty", err)
		}
		return nil
	default:
		return ErrUnknownOpp
	}
}

// 2	Удержание НДФЛ по купонам.
// 8    Удержание налога по дивидендам.
func (p *ReportPositions) ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(
	operation domain.OperationWithoutCustomTypes,
) error {
	return p.distributePayment(
		operation.Payment,
		func(pos *PositionByFIFO) *float64 {
			return &pos.PaidTax
		})
}

// 10	Частичное погашение облигаций.
func (p *ReportPositions) ProcessPartialRedemptionOfBonds(
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	return p.distributePayment(
		operation.Payment,
		func(pos *PositionByFIFO) *float64 {
			return &pos.PartialEarlyRepayment
		})
}

// 15	Покупка ЦБ.
// 16	Покупка ЦБ с карты.
// 57   Перевод ценных бумаг с ИИС на Брокерский счет
func (p *ReportPositions) ProcessPurchaseOfSecurities(
	operation domain.OperationWithoutCustomTypes,
	bondIdentifiers domain.BondIdentIdentifiers,
	lastPrice domain.LastPrice,
	vunitRate domain.Rate,
) {
	const op = "report.processPurchaseOfSecurities"

	// TODO: при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := NewPositionByFIFOFromOperation(operation)

	// Для Евротранса исключение
	// TODO: Это исключение только для одного аккаунта должно работать. Потенциальный баг
	if operation.InstrumentUid == EuroTransInstrumentUID && operation.Type == TransferOfSecuritiesFromAnotherDepository {
		position.BuyPrice = EuroTransBuyCost
	}

	position.ApplyBondMetadata(bondIdentifiers) // Применяем BondIdentifiers из ReportLine.Bond

	nominal := CalculateNominal(bondIdentifiers.Nominal, position.Replaced, vunitRate) // Рассчитываем номинал с учетом курс

	sellPrice := CalculateSellPrice(nominal, lastPrice) // рассчитываем цену продажи

	// Применяем полученные значения к positionbyFIFO
	position.Nominal = nominal
	position.SellPrice = sellPrice

	// Добавляем позицию в ReportPosition.CurrentPositions(открытые позиции на счете)
	p.CurrentPositions = append(p.CurrentPositions, position)
	// Изменяем общее количество бумаг на счете. ReportPosition.Quantity
	p.Quantity += operation.QuantityDone
}

// 21	Выплата дивидендов.
func (p *ReportPositions) ProcessPaymentOfDividends(
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	return p.distributePayment(
		operation.Payment,
		func(pos *PositionByFIFO) *float64 {
			return &pos.TotalDividend
		})
}

// 23 Выплата купонов.
func (p *ReportPositions) ProcessPaymentOfCoupons(
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	return p.distributePayment(
		operation.Payment,
		func(pos *PositionByFIFO) *float64 {
			return &pos.TotalCoupon
		})
}

// 47	Гербовый сбор.
func (p *ReportPositions) ProcessStampDuty(
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	return p.distributePayment(
		operation.Payment,
		func(pos *PositionByFIFO) *float64 {
			return &pos.TotalComission
		})
}

// 22	Продажа ЦБ.
func (p *ReportPositions) ProcessSellOfSecurities(
	operation *domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "report.processSellOfSecurities"

	// Уменьшаем общее количество бумаг на количество проданных
	p.Quantity -= operation.QuantityDone
	// TODO: Переписать ПОЗЖЕ Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления в которых Кол-во проданных
	// бумаг больше кол-ва бумаг в текущей позиции
	var deleteCount int
	// Лейбл end нужен для того чтобы можно было зделать выход из цикла внутри конструкции switch-case
end:
	// Итерируемся по текущим  позициям
	for i := range p.CurrentPositions {
		// Создаем переменную, которая содержит в себе указатель на позицию PositionByFIFO
		currPosition := &p.CurrentPositions[i]
		// Создаем переменную количества бумаг в текущей позиции
		currentQuantity := currPosition.Quantity
		// Создаем переменную количества бумаг в операции продажи
		sellQuantity := operation.QuantityDone
		// Три варината возможно:
		// 1. В текущей позиции больше бумаг, чем в операции продажи
		// 2. Количество равно
		// 3. В операции продажи кол-во бумаг больше , чем в текущей позиции
		switch {
		// 1. В текущей позиции больше бумаг, чем в операции продажи
		case currentQuantity > sellQuantity:
			err := currPosition.isCurrentQuantityGreaterThanSellQuantity(operation.QuantityDone)
			if err != nil {
				return e.WrapIfErr("failed to isCurrentQuantityGreaterThanSellQuantity", err)
			}
			// Прерываем цикл
			break end
		case currPosition.Quantity == operation.QuantityDone:
			p.isEqualCurrentQuantityAndSellQuantity()
			break end
		case currentQuantity < sellQuantity:
			proportion := currentQuantity / sellQuantity
			// Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления
			deleteCount += 1
			operation.ApplyValuesIfCurrentQuantityLessThanSellQuantity(proportion, currentQuantity)
		}

	}
	// удаляем закрытые позиции из среза текущих позиций
	p.CurrentPositions = p.CurrentPositions[deleteCount:]
	return nil
}

func (p *ReportPositions) isEqualCurrentQuantityAndSellQuantity() {
	// Просто закрываем позицию , ведь сколько куплено, столько и продано
	p.CurrentPositions = p.CurrentPositions[1:]
}

func (p *ReportPositions) distributePayment(
	payment float64,
	targetField func(*PositionByFIFO) *float64,
) error {
	if p.Quantity == 0 {
		return ErrZeroQuantity
	}

	for i := range p.CurrentPositions {
		pos := &p.CurrentPositions[i]
		proportion := pos.Quantity / p.Quantity
		*targetField(pos) += payment * proportion
	}

	return nil
}
