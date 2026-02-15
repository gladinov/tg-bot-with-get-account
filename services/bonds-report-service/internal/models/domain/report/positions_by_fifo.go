package report

import (
	"bonds-report-service/internal/models/domain"
	"time"
)

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

func NewPositionByFIFOFromOperation(op domain.OperationWithoutCustomTypes) PositionByFIFO {
	return PositionByFIFO{
		Name:           op.Name,
		BuyDate:        op.Date,
		Figi:           op.Figi,
		Quantity:       op.QuantityDone,
		InstrumentType: op.InstrumentType,
		InstrumentUid:  op.InstrumentUid,
		BuyPrice:       op.Price,
		Currency:       op.Currency,
		BuyAccruedInt:  op.AccruedInt, // НКД при покупке
		TotalComission: op.Commission,
	}
}

func (P *PositionByFIFO) ApplyBondMetadata(
	bond domain.BondIdentIdentifiers,
) {
	P.Ticker = bond.Ticker
	P.ClassCode = bond.ClassCode
	P.Replaced = bond.Replaced
	P.CurrencyIfReplaced = bond.NominalCurrency
}

func (P *PositionByFIFO) isCurrentQuantityGreaterThanSellQuantity(
	sellQuantity float64,
) error {
	// Создаем переменную количества бумаг в текущей позиции
	currentQuantity := P.Quantity
	// Создаем переменную количества бумаг в операции продажи
	// Создаем пересменную пропорции
	var proportion float64
	// Проверяем делитель на ноль
	if currentQuantity == 0 {
		return ErrZeroQuantity
	}
	// Количество проданных бумаг, меньше кол-ва бумаг в позиции.
	//  Получаем пропорцию Проданные бумаги/Текущие бумаги
	proportion = sellQuantity / currentQuantity
	// Отнимаем кол-во проданных бумаг от количества бумаг в текущей позиции
	P.Quantity -= sellQuantity
	// Изменяем значения текущей позиции, умножая на остаток от пропорции
	P.TotalComission = P.TotalComission * (1 - proportion)
	P.PaidTax = P.PaidTax * (1 - proportion)
	P.BuyAccruedInt = P.BuyAccruedInt * (1 - proportion)
	return nil
}

