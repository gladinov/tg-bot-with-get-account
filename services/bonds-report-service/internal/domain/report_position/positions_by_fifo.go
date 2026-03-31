package report

import (
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils"
	"math"
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

func (p *PositionByFIFO) ApplyBondMetadata(
	bond domain.BondIdentIdentifiers,
) {
	p.Ticker = bond.Ticker
	p.ClassCode = bond.ClassCode
	p.Replaced = bond.Replaced
	p.CurrencyIfReplaced = bond.NominalCurrency
}

func (p *PositionByFIFO) isCurrentQuantityGreaterThanSellQuantity(
	sellQuantity float64,
) error {
	// Создаем переменную количества бумаг в текущей позиции
	currentQuantity := p.Quantity
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
	p.Quantity -= sellQuantity
	// Изменяем значения текущей позиции, умножая на остаток от пропорции
	p.TotalComission = p.TotalComission * (1 - proportion)
	p.PaidTax = p.PaidTax * (1 - proportion)
	p.BuyAccruedInt = p.BuyAccruedInt * (1 - proportion)
	return nil
}

func (p *PositionByFIFO) GetProfit(profit float64) (_ float64, err error) {
	totalInvest := p.BuyPrice * p.Quantity
	if totalInvest == 0 {
		return 0, ErrZeroDivision
	}

	profitInPercentage := utils.RoundFloat((profit/(totalInvest))*100, 2)
	return profitInPercentage, nil
}

func (p *PositionByFIFO) GetAnnualizedReturnInPercentage(netProfit float64, sellDate time.Time) (_ float64, err error) {
	totalInvested := p.BuyPrice * p.Quantity
	if totalInvested == 0 {
		return 0, ErrZeroDivision
	}
	if p.BuyDate.After(sellDate) {
		return 0, ErrInvalidDate
	}

	totalReturn := netProfit / totalInvested

	buyDate := p.BuyDate

	timeDurationInDays := sellDate.Sub(buyDate).Hours() / float64(hoursInDay)
	// Если покупка и продажа были совершены в один день, то берем минимум один день
	timeDurationInYears := math.Max(1, timeDurationInDays) / float64(daysInYear)

	annualizedReturn := math.Pow((1+totalReturn), (1/timeDurationInYears)) - 1
	annualizedReturnInPercentage := annualizedReturn * 100
	annualizedReturnInPercentageRound := utils.RoundFloat(annualizedReturnInPercentage, 2)

	return annualizedReturnInPercentageRound, nil
}

// Доход по позиции до вычета налога
func (p *PositionByFIFO) GetProfitBeforeTax() float64 {
	// Для незакрытых позиций необходимо посчитать еще потенциальную комиссию при продаже
	quantity := p.Quantity
	buySellDifference := (p.SellPrice-p.BuyPrice)*quantity + p.SellAccruedInt - p.BuyAccruedInt
	cashFlow := p.TotalCoupon + p.TotalDividend
	positionProfit := buySellDifference + cashFlow + p.TotalComission + p.PartialEarlyRepayment
	return positionProfit
}

// Расход полного налога по закрытой позиции
func (p *PositionByFIFO) GetTotalTaxFromPosition(profit float64) float64 {
	// Рассчитываем срок владения
	buyDate := p.BuyDate
	sellDate := p.SellDate

	// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
	if profit < 0 || isHoldingPeriodMoreThanThreeYears(buyDate, sellDate) {
		return 0
	}
	return profit * baseTaxRate
}

func isHoldingPeriodMoreThanThreeYears(buyDate time.Time, sellDate time.Time) bool {
	threeYearAfterBuyDate := buyDate.AddDate(3, 0, 0)
	if threeYearAfterBuyDate.After(sellDate) {
		return false
	}
	return true
}
