//go:build unit

package report

import (
	"bonds-report-service/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDistributePayment(t *testing.T) {
	t.Run("успешное распределение платежа по позициям", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0, TotalDividend: 100.0},
				{Quantity: 50.0, TotalDividend: 200.0},
				{Quantity: 20.0, TotalDividend: 50.0},
			},
		}

		err := report.distributePayment(500.0, func(p *PositionByFIFO) *float64 {
			return &p.TotalDividend
		})

		require.NoError(t, err)
		require.Equal(t, 250.0, report.CurrentPositions[0].TotalDividend) // 100 + 500*0.3
		require.Equal(t, 450.0, report.CurrentPositions[1].TotalDividend) // 200 + 500*0.5
		require.Equal(t, 150.0, report.CurrentPositions[2].TotalDividend) // 50 + 500*0.2
	})

	t.Run("нулевое общее количество - ошибка", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0},
			},
		}

		err := report.distributePayment(500.0, func(p *PositionByFIFO) *float64 {
			return &p.TotalDividend
		})

		require.ErrorIs(t, err, ErrZeroQuantity)
	})

	t.Run("пустой список позиций - не паникует", func(t *testing.T) {
		report := &ReportPositions{
			Quantity:         100.0,
			CurrentPositions: []PositionByFIFO{},
		}

		err := report.distributePayment(500.0, func(p *PositionByFIFO) *float64 {
			return &p.TotalDividend
		})

		require.NoError(t, err)
	})

	t.Run("дробные количества - точность вычислений", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.5,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.3, TotalDividend: 100.0},
				{Quantity: 70.2, TotalDividend: 200.0},
			},
		}

		err := report.distributePayment(500.0, func(p *PositionByFIFO) *float64 {
			return &p.TotalDividend
		})

		require.NoError(t, err)

		// 30.3/100.5 ≈ 0.3015 * 500 = 150.75
		require.InDelta(t, 250.75, report.CurrentPositions[0].TotalDividend, 0.05)

		// 70.2/100.5 ≈ 0.6985 * 500 = 349.25
		require.InDelta(t, 549.25, report.CurrentPositions[1].TotalDividend, 0.05)
	})

	t.Run("нулевой платеж - поля не меняются", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0, TotalDividend: 100.0},
				{Quantity: 70.0, TotalDividend: 200.0},
			},
		}

		original := *report

		err := report.distributePayment(0.0, func(p *PositionByFIFO) *float64 {
			return &p.TotalDividend
		})

		require.NoError(t, err)
		require.Equal(t, original.CurrentPositions[0].TotalDividend, report.CurrentPositions[0].TotalDividend)
		require.Equal(t, original.CurrentPositions[1].TotalDividend, report.CurrentPositions[1].TotalDividend)
	})
}

func TestPaymentProcessors(t *testing.T) {
	// Общие данные для всех тестов
	operation := domain.OperationWithoutCustomTypes{Payment: 500.0}

	t.Run("ProcessPaymentOfDividends обновляет поле TotalDividend", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0},
				{Quantity: 70.0},
			},
		}

		report.CurrentPositions[0].TotalDividend = 100.0
		report.CurrentPositions[1].TotalDividend = 200.0

		err := report.ProcessPaymentOfDividends(operation)
		require.NoError(t, err)

		// Проверяем, что обновилось только поле TotalDividend
		pos1 := report.CurrentPositions[0]
		require.Equal(t, 250.0, pos1.TotalDividend) // 100 + 500*0.3
		require.Zero(t, pos1.TotalCoupon)
		require.Zero(t, pos1.TotalComission)
		require.Zero(t, pos1.PaidTax)
		require.Zero(t, pos1.PartialEarlyRepayment)

		pos2 := report.CurrentPositions[1]
		require.Equal(t, 550.0, pos2.TotalDividend) // 200 + 500*0.7
		require.Zero(t, pos2.TotalCoupon)
		require.Zero(t, pos2.TotalComission)
		require.Zero(t, pos2.PaidTax)
		require.Zero(t, pos2.PartialEarlyRepayment)
	})

	t.Run("ProcessPaymentOfCoupons обновляет поле TotalCoupon", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0},
				{Quantity: 70.0},
			},
		}

		report.CurrentPositions[0].TotalCoupon = 100.0
		report.CurrentPositions[1].TotalCoupon = 200.0

		err := report.ProcessPaymentOfCoupons(operation)
		require.NoError(t, err)

		pos1 := report.CurrentPositions[0]
		require.Equal(t, 250.0, pos1.TotalCoupon) // 100 + 500*0.3
		require.Zero(t, pos1.TotalDividend)
		require.Zero(t, pos1.TotalComission)
		require.Zero(t, pos1.PaidTax)
		require.Zero(t, pos1.PartialEarlyRepayment)

		pos2 := report.CurrentPositions[1]
		require.Equal(t, 550.0, pos2.TotalCoupon) // 200 + 500*0.7
		require.Zero(t, pos2.TotalDividend)
		require.Zero(t, pos2.TotalComission)
		require.Zero(t, pos2.PaidTax)
		require.Zero(t, pos2.PartialEarlyRepayment)
	})

	t.Run("ProcessStampDuty обновляет поле TotalComission", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0},
				{Quantity: 70.0},
			},
		}

		report.CurrentPositions[0].TotalComission = 100.0
		report.CurrentPositions[1].TotalComission = 200.0

		err := report.ProcessStampDuty(operation)
		require.NoError(t, err)

		pos1 := report.CurrentPositions[0]
		require.Equal(t, 250.0, pos1.TotalComission) // 100 + 500*0.3
		require.Zero(t, pos1.TotalDividend)
		require.Zero(t, pos1.TotalCoupon)
		require.Zero(t, pos1.PaidTax)
		require.Zero(t, pos1.PartialEarlyRepayment)

		pos2 := report.CurrentPositions[1]
		require.Equal(t, 550.0, pos2.TotalComission) // 200 + 500*0.7
		require.Zero(t, pos2.TotalDividend)
		require.Zero(t, pos2.TotalCoupon)
		require.Zero(t, pos2.PaidTax)
		require.Zero(t, pos2.PartialEarlyRepayment)
	})

	t.Run("ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends обновляет поле PaidTax", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0},
				{Quantity: 70.0},
			},
		}

		report.CurrentPositions[0].PaidTax = 100.0
		report.CurrentPositions[1].PaidTax = 200.0

		err := report.ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(operation)
		require.NoError(t, err)

		pos1 := report.CurrentPositions[0]
		require.Equal(t, 250.0, pos1.PaidTax) // 100 + 500*0.3
		require.Zero(t, pos1.TotalDividend)
		require.Zero(t, pos1.TotalCoupon)
		require.Zero(t, pos1.TotalComission)
		require.Zero(t, pos1.PartialEarlyRepayment)

		pos2 := report.CurrentPositions[1]
		require.Equal(t, 550.0, pos2.PaidTax) // 200 + 500*0.7
		require.Zero(t, pos2.TotalDividend)
		require.Zero(t, pos2.TotalCoupon)
		require.Zero(t, pos2.TotalComission)
		require.Zero(t, pos2.PartialEarlyRepayment)
	})

	t.Run("ProcessPartialRedemptionOfBonds обновляет поле PartialEarlyRepayment", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0},
				{Quantity: 70.0},
			},
		}

		report.CurrentPositions[0].PartialEarlyRepayment = 100.0
		report.CurrentPositions[1].PartialEarlyRepayment = 200.0

		err := report.ProcessPartialRedemptionOfBonds(operation)
		require.NoError(t, err)

		pos1 := report.CurrentPositions[0]
		require.Equal(t, 250.0, pos1.PartialEarlyRepayment) // 100 + 500*0.3
		require.Zero(t, pos1.TotalDividend)
		require.Zero(t, pos1.TotalCoupon)
		require.Zero(t, pos1.TotalComission)
		require.Zero(t, pos1.PaidTax)

		pos2 := report.CurrentPositions[1]
		require.Equal(t, 550.0, pos2.PartialEarlyRepayment) // 200 + 500*0.7
		require.Zero(t, pos2.TotalDividend)
		require.Zero(t, pos2.TotalCoupon)
		require.Zero(t, pos2.TotalComission)
		require.Zero(t, pos2.PaidTax)
	})

	t.Run("все процессоры возвращают ошибку при нулевом Quantity", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 0,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.0},
			},
		}

		// Проверяем каждый процессор отдельно
		err := report.ProcessPaymentOfDividends(operation)
		require.ErrorIs(t, err, ErrZeroQuantity)

		err = report.ProcessPaymentOfCoupons(operation)
		require.ErrorIs(t, err, ErrZeroQuantity)

		err = report.ProcessStampDuty(operation)
		require.ErrorIs(t, err, ErrZeroQuantity)

		err = report.ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(operation)
		require.ErrorIs(t, err, ErrZeroQuantity)

		err = report.ProcessPartialRedemptionOfBonds(operation)
		require.ErrorIs(t, err, ErrZeroQuantity)
	})

	t.Run("процессоры работают с дробными количествами", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.5,
			CurrentPositions: []PositionByFIFO{
				{Quantity: 30.3, TotalDividend: 100.0},
				{Quantity: 70.2, TotalDividend: 200.0},
			},
		}

		err := report.ProcessPaymentOfDividends(operation)
		require.NoError(t, err)

		// 30.3/100.5 ≈ 0.3015 * 500 = 150.75
		require.InDelta(t, 250.75, report.CurrentPositions[0].TotalDividend, 0.05)

		// 70.2/100.5 ≈ 0.6985 * 500 = 349.25
		require.InDelta(t, 549.25, report.CurrentPositions[1].TotalDividend, 0.05)
	})
}

func TestPaymentProcessors_Errors(t *testing.T) {
	report := &ReportPositions{
		Quantity: 0,
		CurrentPositions: []PositionByFIFO{
			{Quantity: 30.0},
		},
	}

	operation := domain.OperationWithoutCustomTypes{Payment: 500.0}

	// Проверяем все обертки на одной ошибке
	processors := []func(*ReportPositions, domain.OperationWithoutCustomTypes) error{
		(*ReportPositions).ProcessPaymentOfDividends,
		(*ReportPositions).ProcessPaymentOfCoupons,
		(*ReportPositions).ProcessStampDuty,
		(*ReportPositions).ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends,
		(*ReportPositions).ProcessPartialRedemptionOfBonds,
	}

	for _, proc := range processors {
		err := proc(report, operation)
		require.ErrorIs(t, err, ErrZeroQuantity)
	}
}

func TestProcessPurchaseOfSecurities(t *testing.T) {
	// Базовые тестовые данные
	bondIdentifiers := domain.BondIdentIdentifiers{
		Ticker:          "SU29006RMFS5",
		ClassCode:       "TQOB",
		Replaced:        true,
		NominalCurrency: "USD",
		Nominal:         domain.MoneyValue{Units: 1000, Nano: 0},
	}

	lastPrice := domain.LastPrice{
		LastPrice: domain.Quotation{Units: 98, Nano: 500000000}, // 98.5
	}

	vunitRate := domain.Rate{
		IsoCurrencyName: "USD",
		Vunit_Rate: domain.NullFloat64{
			Value:  1.2, // курс 1.2
			IsSet:  true,
			IsNull: false,
		},
	}

	t.Run("добавляет новую позицию при покупке", func(t *testing.T) {
		report := &ReportPositions{
			Quantity:         100.0,
			CurrentPositions: []PositionByFIFO{},
		}

		operation := domain.OperationWithoutCustomTypes{
			Name:           "Покупка облигаций",
			QuantityDone:   10.0,
			Price:          1050.5,
			Payment:        10505.0,
			Date:           time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			Commission:     100.0,
			AccruedInt:     50.0,
			InstrumentType: "Bond",
			Figi:           "BBG000000001",
			InstrumentUid:  "UID123",
			Currency:       "RUB",
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, vunitRate)

		require.Len(t, report.CurrentPositions, 1)

		pos := report.CurrentPositions[0]

		// Проверяем поля из NewPositionByFIFOFromOperation
		require.Equal(t, operation.Name, pos.Name)
		require.Equal(t, operation.Date, pos.BuyDate)
		require.Equal(t, operation.Figi, pos.Figi)
		require.Equal(t, operation.QuantityDone, pos.Quantity)
		require.Equal(t, operation.InstrumentType, pos.InstrumentType)
		require.Equal(t, operation.InstrumentUid, pos.InstrumentUid)
		require.Equal(t, operation.Price, pos.BuyPrice)
		require.Equal(t, operation.Currency, pos.Currency)
		require.Equal(t, operation.AccruedInt, pos.BuyAccruedInt)
		require.Equal(t, operation.Commission, pos.TotalComission)

		// Проверяем ApplyBondMetadata
		require.Equal(t, bondIdentifiers.Ticker, pos.Ticker)
		require.Equal(t, bondIdentifiers.ClassCode, pos.ClassCode)
		require.Equal(t, bondIdentifiers.Replaced, pos.Replaced)
		require.Equal(t, bondIdentifiers.NominalCurrency, pos.CurrencyIfReplaced)

		// Проверяем CalculateNominal
		expectedNominal := 1200.0 // 1000 * 1.2 (replaced=true, rate.IsSet=true)
		require.Equal(t, expectedNominal, pos.Nominal)

		// Проверяем CalculateSellPrice
		// lastPrice = 98.5, nominal = 1200
		// 98.5/100 * 1200 = 1182
		expectedSellPrice := 1182.0
		require.Equal(t, expectedSellPrice, pos.SellPrice)

		// Поля продажи не должны быть заполнены
		require.True(t, pos.SellDate.IsZero())
		require.Zero(t, pos.SellAccruedInt)
		require.Zero(t, pos.SellPayment)

		// Общее количество увеличилось
		require.Equal(t, 110.0, report.Quantity)
	})

	t.Run("покупка без замены валюты", func(t *testing.T) {
		report := &ReportPositions{}

		bondWithoutReplacement := domain.BondIdentIdentifiers{
			Ticker:          "SU29006RMFS5",
			ClassCode:       "TQOB",
			Replaced:        false, // replaced = false
			NominalCurrency: "RUB",
			Nominal:         domain.MoneyValue{Units: 1000, Nano: 0},
		}

		operation := domain.OperationWithoutCustomTypes{
			QuantityDone: 10.0,
			Price:        1000.0,
		}

		report.ProcessPurchaseOfSecurities(operation, bondWithoutReplacement, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// CalculateNominal: replaced=false -> должен вернуть base
		require.Equal(t, 1000.0, pos.Nominal)

		// Replaced должно быть false из метаданных
		require.False(t, pos.Replaced)
	})

	t.Run("покупка с rate не установлен", func(t *testing.T) {
		report := &ReportPositions{}

		rateNotSet := domain.Rate{
			Vunit_Rate: domain.NullFloat64{
				Value:  1.2,
				IsSet:  false, // не установлен
				IsNull: false,
			},
		}

		operation := domain.OperationWithoutCustomTypes{
			QuantityDone: 10.0,
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, rateNotSet)

		pos := report.CurrentPositions[0]

		// CalculateNominal: replaced=true, но IsSet=false -> должен вернуть base
		require.Equal(t, 1000.0, pos.Nominal)
	})

	t.Run("покупка с rate is null", func(t *testing.T) {
		report := &ReportPositions{}

		rateIsNull := domain.Rate{
			Vunit_Rate: domain.NullFloat64{
				Value:  1.2,
				IsSet:  true,
				IsNull: true, // is null
			},
		}

		operation := domain.OperationWithoutCustomTypes{
			QuantityDone: 10.0,
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, rateIsNull)

		pos := report.CurrentPositions[0]

		// CalculateNominal: replaced=true, IsNull=true -> должен вернуть base
		require.Equal(t, 1000.0, pos.Nominal)
	})

	t.Run("разные цены для CalculateSellPrice", func(t *testing.T) {
		operation := domain.OperationWithoutCustomTypes{
			QuantityDone: 10.0,
		}

		testCases := []struct {
			name     string
			units    int64
			nano     int32
			nominal  float64
			expected float64
		}{
			{
				name:     "цена 100%",
				units:    100,
				nano:     0,
				nominal:  1000,
				expected: 1000.0,
			},
			{
				name:     "цена 98.5%",
				units:    98,
				nano:     500000000,
				nominal:  1000,
				expected: 985.0,
			},
			{
				name:     "цена 105.75%",
				units:    105,
				nano:     750000000,
				nominal:  1000,
				expected: 1057.5,
			},
			{
				name:     "с другим номиналом",
				units:    98,
				nano:     500000000,
				nominal:  1200,
				expected: 1182.0,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Создаем новый отчет для каждого кейса
				r := &ReportPositions{}

				price := domain.LastPrice{
					LastPrice: domain.Quotation{
						Units: tc.units,
						Nano:  tc.nano,
					},
				}

				// Для теста создаем bond с нужным номиналом
				bond := domain.BondIdentIdentifiers{
					Replaced: false,
					Nominal:  domain.MoneyValue{Units: int64(tc.nominal), Nano: 0},
				}

				r.ProcessPurchaseOfSecurities(operation, bond, price, vunitRate)

				pos := r.CurrentPositions[0]
				require.InDelta(t, tc.expected, pos.SellPrice, 0.01)
			})
		}
	})

	t.Run("дробный номинал в MoneyValue", func(t *testing.T) {
		report := &ReportPositions{}

		bondWithFractional := domain.BondIdentIdentifiers{
			Replaced: false,
			Nominal:  domain.MoneyValue{Units: 1234, Nano: 560000000}, // 1234.56
		}

		operation := domain.OperationWithoutCustomTypes{
			QuantityDone: 10.0,
		}

		report.ProcessPurchaseOfSecurities(operation, bondWithFractional, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// CalculateNominal должен правильно сконвертировать MoneyValue в float64
		require.InDelta(t, 1234.56, pos.Nominal, 0.0001)
	})

	t.Run("комбинация всех функций", func(t *testing.T) {
		report := &ReportPositions{}

		operation := domain.OperationWithoutCustomTypes{
			Name:           "Тестовая операция",
			QuantityDone:   5.0,
			Price:          950.0,
			Date:           time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			Commission:     25.0,
			AccruedInt:     15.0,
			InstrumentType: "Bond",
			Figi:           "TEST_FIGI",
			Currency:       "USD",
		}

		customBond := domain.BondIdentIdentifiers{
			Ticker:          "TEST_TICKER",
			ClassCode:       "TEST_CLASS",
			Replaced:        true,
			NominalCurrency: "USD",
			Nominal:         domain.MoneyValue{Units: 500, Nano: 0},
		}

		customRate := domain.Rate{
			Vunit_Rate: domain.NullFloat64{
				Value:  2.0, // курс 2.0
				IsSet:  true,
				IsNull: false,
			},
		}

		customPrice := domain.LastPrice{
			LastPrice: domain.Quotation{Units: 101, Nano: 500000000}, // 101.5
		}

		report.ProcessPurchaseOfSecurities(operation, customBond, customPrice, customRate)

		pos := report.CurrentPositions[0]

		// 1. NewPositionByFIFOFromOperation
		require.Equal(t, "Тестовая операция", pos.Name)
		require.Equal(t, 5.0, pos.Quantity)
		require.Equal(t, 950.0, pos.BuyPrice)
		require.Equal(t, "TEST_FIGI", pos.Figi)
		require.Equal(t, 25.0, pos.TotalComission)
		require.Equal(t, 15.0, pos.BuyAccruedInt)

		// 2. ApplyBondMetadata
		require.Equal(t, "TEST_TICKER", pos.Ticker)
		require.Equal(t, "TEST_CLASS", pos.ClassCode)
		require.True(t, pos.Replaced)
		require.Equal(t, "USD", pos.CurrencyIfReplaced)

		// 3. CalculateNominal
		expectedNominal := 1000.0 // 500 * 2.0
		require.Equal(t, expectedNominal, pos.Nominal)

		// 4. CalculateSellPrice
		expectedSellPrice := 1015.0 // 101.5/100 * 1000 = 1015
		require.Equal(t, expectedSellPrice, pos.SellPrice)

		// Проверяем общее количество
		require.Equal(t, 5.0, report.Quantity)
	})
	t.Run("евротранс с правильным UID и типом операции - переопределяет BuyPrice", func(t *testing.T) {
		report := &ReportPositions{}

		operation := domain.OperationWithoutCustomTypes{
			Name:           "Перевод еврооблигаций",
			QuantityDone:   5.0,
			Price:          950.0, // оригинальная цена из операции
			Date:           time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			InstrumentUid:  EuroTransInstrumentUID,                    // специальный UID
			Type:           TransferOfSecuritiesFromAnotherDepository, // специальный тип
			InstrumentType: "Bond",
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// BuyPrice должен быть переопределен на EuroTransBuyCost
		require.Equal(t, EuroTransBuyCost, pos.BuyPrice)
		// Не должен использовать оригинальную цену из операции
		require.NotEqual(t, operation.Price, pos.BuyPrice)
	})

	t.Run("евротранс с правильным UID но неправильным типом - не переопределяет", func(t *testing.T) {
		report := &ReportPositions{}

		operation := domain.OperationWithoutCustomTypes{
			Name:           "Перевод с другим типом",
			QuantityDone:   5.0,
			Price:          950.0,
			InstrumentUid:  EuroTransInstrumentUID, // специальный UID
			Type:           15,                     // НЕ тот тип (обычная покупка)
			InstrumentType: "Bond",
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// BuyPrice должен быть из операции, не переопределен
		require.Equal(t, operation.Price, pos.BuyPrice)
		require.NotEqual(t, EuroTransBuyCost, pos.BuyPrice)
	})

	t.Run("евротранс с правильным типом но неправильным UID - не переопределяет", func(t *testing.T) {
		report := &ReportPositions{}

		operation := domain.OperationWithoutCustomTypes{
			Name:           "Перевод с другим UID",
			QuantityDone:   5.0,
			Price:          950.0,
			InstrumentUid:  "some-other-uid",                          // НЕ специальный UID
			Type:           TransferOfSecuritiesFromAnotherDepository, // специальный тип
			InstrumentType: "Bond",
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// BuyPrice должен быть из операции, не переопределен
		require.Equal(t, operation.Price, pos.BuyPrice)
	})

	t.Run("несколько операций - евротранс влияет только на свою позицию", func(t *testing.T) {
		report := &ReportPositions{}

		// Сначала обычная покупка
		op1 := domain.OperationWithoutCustomTypes{
			Name:          "Обычная покупка",
			QuantityDone:  10.0,
			Price:         1000.0,
			InstrumentUid: "regular-uid-1",
			Type:          15,
		}
		report.ProcessPurchaseOfSecurities(op1, bondIdentifiers, lastPrice, vunitRate)

		// Потом евротранс
		op2 := domain.OperationWithoutCustomTypes{
			Name:          "Евротранс",
			QuantityDone:  5.0,
			Price:         900.0,
			InstrumentUid: EuroTransInstrumentUID,
			Type:          TransferOfSecuritiesFromAnotherDepository,
		}
		report.ProcessPurchaseOfSecurities(op2, bondIdentifiers, lastPrice, vunitRate)

		// Потом еще одна обычная покупка
		op3 := domain.OperationWithoutCustomTypes{
			Name:          "Еще покупка",
			QuantityDone:  3.0,
			Price:         1100.0,
			InstrumentUid: "regular-uid-2",
			Type:          15,
		}
		report.ProcessPurchaseOfSecurities(op3, bondIdentifiers, lastPrice, vunitRate)

		require.Len(t, report.CurrentPositions, 3)

		// Первая позиция - обычная
		require.Equal(t, 1000.0, report.CurrentPositions[0].BuyPrice)
		require.Equal(t, "regular-uid-1", report.CurrentPositions[0].InstrumentUid)

		// Вторая позиция - евротранс (должна иметь специальную цену)
		require.Equal(t, EuroTransBuyCost, report.CurrentPositions[1].BuyPrice)
		require.Equal(t, EuroTransInstrumentUID, report.CurrentPositions[1].InstrumentUid)

		// Третья позиция - обычная
		require.Equal(t, 1100.0, report.CurrentPositions[2].BuyPrice)
		require.Equal(t, "regular-uid-2", report.CurrentPositions[2].InstrumentUid)

		// Общее количество
		require.Equal(t, 18.0, report.Quantity)
	})

	t.Run("евротранс с другими типами инструментов", func(t *testing.T) {
		report := &ReportPositions{}

		// Евротранс для акций
		op := domain.OperationWithoutCustomTypes{
			Name:           "Евротранс акций",
			QuantityDone:   5.0,
			Price:          950.0,
			InstrumentUid:  EuroTransInstrumentUID,
			Type:           TransferOfSecuritiesFromAnotherDepository,
			InstrumentType: "Share", // акции, а не облигации
		}

		report.ProcessPurchaseOfSecurities(op, bondIdentifiers, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// Проверяем, что специальная логика работает для всех типов инструментов
		require.Equal(t, EuroTransBuyCost, pos.BuyPrice)
		require.Equal(t, "Share", pos.InstrumentType)
	})

	t.Run("евротранс с нулевым количеством", func(t *testing.T) {
		report := &ReportPositions{
			Quantity: 100.0,
		}

		operation := domain.OperationWithoutCustomTypes{
			Name:          "Евротранс с нулем",
			QuantityDone:  0.0,
			Price:         950.0,
			InstrumentUid: EuroTransInstrumentUID,
			Type:          TransferOfSecuritiesFromAnotherDepository,
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// BuyPrice все равно должен быть переопределен
		require.Equal(t, EuroTransBuyCost, pos.BuyPrice)
		require.Equal(t, 0.0, pos.Quantity)
		require.Equal(t, 100.0, report.Quantity) // общее количество не изменилось
	})

	t.Run("все поля для евротранс заполняются корректно", func(t *testing.T) {
		report := &ReportPositions{}

		operation := domain.OperationWithoutCustomTypes{
			Name:           "Евротранс полный тест",
			QuantityDone:   7.0,
			Price:          950.0, // будет переопределено
			Date:           time.Date(2024, 1, 15, 15, 30, 0, 0, time.UTC),
			Commission:     50.0,
			AccruedInt:     25.0,
			InstrumentType: "Bond",
			Figi:           "EURO_FIGI",
			InstrumentUid:  EuroTransInstrumentUID,
			Type:           TransferOfSecuritiesFromAnotherDepository,
			Currency:       "USD",
			Payment:        7000.0,
		}

		report.ProcessPurchaseOfSecurities(operation, bondIdentifiers, lastPrice, vunitRate)

		pos := report.CurrentPositions[0]

		// Проверяем все поля
		require.Equal(t, operation.Name, pos.Name)
		require.Equal(t, operation.Date, pos.BuyDate)
		require.Equal(t, operation.Figi, pos.Figi)
		require.Equal(t, operation.QuantityDone, pos.Quantity)
		require.Equal(t, operation.InstrumentType, pos.InstrumentType)
		require.Equal(t, operation.InstrumentUid, pos.InstrumentUid)
		require.Equal(t, operation.Currency, pos.Currency)
		require.Equal(t, operation.AccruedInt, pos.BuyAccruedInt)
		require.Equal(t, operation.Commission, pos.TotalComission)

		// Специально переопределенное поле
		require.Equal(t, EuroTransBuyCost, pos.BuyPrice)
		require.NotEqual(t, operation.Price, pos.BuyPrice)

		// Метаданные
		require.Equal(t, bondIdentifiers.Ticker, pos.Ticker)
		require.Equal(t, bondIdentifiers.ClassCode, pos.ClassCode)

		// Рассчитанные поля
		require.Equal(t, 1200.0, pos.Nominal)
		require.Equal(t, 1182.0, pos.SellPrice)
	})
}

func TestIsEqualCurrentQuantityAndSellQuantity(t *testing.T) {
	t.Run("удаляет первую позицию", func(t *testing.T) {
		report := &ReportPositions{
			CurrentPositions: []PositionByFIFO{
				{Name: "Позиция 1", Quantity: 10},
				{Name: "Позиция 2", Quantity: 20},
				{Name: "Позиция 3", Quantity: 30},
			},
		}

		report.isEqualCurrentQuantityAndSellQuantity()

		expected := []PositionByFIFO{
			{Name: "Позиция 2", Quantity: 20},
			{Name: "Позиция 3", Quantity: 30},
		}

		require.Equal(t, expected, report.CurrentPositions)
	})
}
