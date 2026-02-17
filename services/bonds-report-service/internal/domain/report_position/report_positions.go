package report

import (
	"bonds-report-service/internal/domain"

	"github.com/gladinov/e"
)

// TODO: Подумать как разбить адекватно эту мешанину из файлов

const (
	WithholdingOfPersonalIncomeTaxOnCoupons        = 2  // 2	Удержание НДФЛ по купонам.
	WithholdingOfPersonalIncomeTaxOnDividends      = 8  // 8    Удержание налога по дивидендам.
	PartialRedemptionOfBonds                       = 10 // 10	Частичное погашение облигаций.
	PurchaseOfSecurities                           = 15 // 15	Покупка ЦБ.
	PurchaseOfSecuritiesWithACard                  = 16 // 16	Покупка ЦБ с карты.
	TransferOfSecuritiesFromAnotherDepository      = 17 // 17	Перевод ценных бумаг из другого депозитария.
	WithhouldingACommissionForTheTransaction       = 19 // 19	Удержание комиссии за операцию.
	PaymentOfDividends                             = 21 // 21	Выплата дивидендов.
	SaleOfSecurities                               = 22 // 22	Продажа ЦБ.
	PaymentOfCoupons                               = 23 // 23 Выплата купонов.
	StampDuty                                      = 47 // 47	Гербовый сбор.
	TransferOfSecuritiesFromIISToABrokerageAccount = 57 // 57   Перевод ценных бумаг с ИИС на Брокерский счет
)

const (
	EuroTransBuyCost       = 240 // Стоимость Евротранса при переводе из другого депозитария
	EuroTransInstrumentUID = "02b2ea14-3c4b-47e8-9548-45a8dbcc8f8a"
	threeYearInHours       = 26304 // Три года в часах
	baseTaxRate            = 0.13  // Налог с продажи ЦБ
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
		// 57   Перевод ценных бумаг с ИИС на Брокерский счет
	case PurchaseOfSecurities, PurchaseOfSecuritiesWithACard, TransferOfSecuritiesFromIISToABrokerageAccount:
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
		// 17	Перевод ценных бумаг из другого депозитария.
	case TransferOfSecuritiesFromAnotherDepository:
		// Проверяем операцию на выполнение.
		// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
		if operation.QuantityDone == 0 {
			return ErrZeroQuantity
		}
		p.ProcessTransferOfSecuritiesFromAnotherDepository(
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
) (err error) {
	const op = "report.processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends"

	// Проверка на наличие бумаг на счете
	if p.Quantity == 0 {
		return ErrZeroQuantity
	} else {
		// Итерируемся по ReportPositions.CurrentPositions
		for i := range p.CurrentPositions {
			// Создаем переменную содержащую в себе ссылку на элемент массива текущих позиций positionByFIFO
			currentPosition := &p.CurrentPositions[i] // TODO: Поинтересоваться у нейронки про экономию ресурсов. Корректно ли я вызвал элемент массива?
			// Рассчитываем пропорцию от общего налога на текущую позицию.
			proportion := currentPosition.Quantity / p.Quantity
			// Плюсуем к уплаченному налогу по текущей позиции выплату по операции * пропорцию
			currentPosition.PaidTax += operation.Payment * proportion
		}
	}
	return nil
}

// 10	Частичное погашение облигаций.
func (p *ReportPositions) ProcessPartialRedemptionOfBonds(
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "report.processPartialRedemptionOfBonds"
	// Описание шагов идентично методу ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends
	if p.Quantity == 0 {
		return ErrZeroQuantity
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
func (p *ReportPositions) ProcessPurchaseOfSecurities(
	operation domain.OperationWithoutCustomTypes,
	bondIdentifiers domain.BondIdentIdentifiers,
	lastPrice domain.LastPrice,
	vunitRate domain.Rate,
) {
	const op = "report.processPurchaseOfSecurities"

	// TODO: при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := NewPositionByFIFOFromOperation(operation) // Создаем PositionByFIFO для операции покупки или перевода бумаги

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

// 17	Перевод ценных бумаг из другого депозитария.
func (p *ReportPositions) ProcessTransferOfSecuritiesFromAnotherDepository(
	operation domain.OperationWithoutCustomTypes,
	bondIdentifiers domain.BondIdentIdentifiers,
	lastPrice domain.LastPrice,
	vunitRate domain.Rate,
) {
	// Описание операций смотри выше
	const op = "report.processTransferOfSecuritiesFromAnotherDepository"

	// TODO: при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := NewPositionByFIFOFromOperation(operation)
	// Для Евротранса исключение // TODO: Это исключение только для одного аккаунта должно работать. Потенциальный баг
	if operation.InstrumentUid == EuroTransInstrumentUID {
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
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "report.processPaymentOfDividends"
	// Описание шагов идентично методу ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends

	if p.Quantity == 0 {
		return ErrZeroQuantity
	} else {
		for i := range p.CurrentPositions {
			currentPosition := &p.CurrentPositions[i]
			proportion := currentPosition.Quantity / p.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			// TODO: Проверить высказывание выше
			currentPosition.TotalDividend += operation.Payment * proportion
		}
	}
	return nil
}

// 23 Выплата купонов.
func (p *ReportPositions) ProcessPaymentOfCoupons(
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "report.processPaymentOfCoupons"
	// Описание шагов идентично методу ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends

	if p.Quantity == 0 {
		return ErrZeroQuantity
	} else {
		for i := range p.CurrentPositions {
			currentPosition := &p.CurrentPositions[i]
			proportion := currentPosition.Quantity / p.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			// TODO: Проверить высказывание выше
			currentPosition.TotalCoupon += operation.Payment * proportion
		}
	}
	return nil
}

// 47	Гербовый сбор.
func (p *ReportPositions) ProcessStampDuty(
	operation domain.OperationWithoutCustomTypes,
) (err error) {
	const op = "report.processStampDuty"
	// Описание шагов идентично методу ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends

	if p.Quantity == 0 {
		return ErrZeroQuantity
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
