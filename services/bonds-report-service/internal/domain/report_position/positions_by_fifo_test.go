//go:build unit

package report

import (
	"bonds-report-service/internal/domain"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestApplyBondMetadata(t *testing.T) {
	t.Run("copies all fields from bond to position", func(t *testing.T) {
		position := &PositionByFIFO{}

		bond := domain.BondIdentIdentifiers{
			Ticker:          "SU29006RMFS5",
			ClassCode:       "TQOB",
			Replaced:        true,
			NominalCurrency: "USD",
		}

		position.ApplyBondMetadata(bond)

		// Проверяем, что все поля скопированы
		if position.Ticker != bond.Ticker {
			t.Errorf("Ticker = %q, want %q", position.Ticker, bond.Ticker)
		}
		if position.ClassCode != bond.ClassCode {
			t.Errorf("ClassCode = %q, want %q", position.ClassCode, bond.ClassCode)
		}
		if position.Replaced != bond.Replaced {
			t.Errorf("Replaced = %v, want %v", position.Replaced, bond.Replaced)
		}
		if position.CurrencyIfReplaced != bond.NominalCurrency {
			t.Errorf("CurrencyIfReplaced = %q, want %q",
				position.CurrencyIfReplaced, bond.NominalCurrency)
		}
	})

	t.Run("overwrites existing values", func(t *testing.T) {
		position := &PositionByFIFO{
			Ticker:             "OLD_TICKER",
			ClassCode:          "OLD_CLASS",
			Replaced:           false,
			CurrencyIfReplaced: "RUB",
		}

		bond := domain.BondIdentIdentifiers{
			Ticker:          "NEW_TICKER",
			ClassCode:       "NEW_CLASS",
			Replaced:        true,
			NominalCurrency: "USD",
		}

		position.ApplyBondMetadata(bond)

		// Проверяем, что старые значения заменены
		if position.Ticker != bond.Ticker {
			t.Errorf("Ticker = %q, want %q", position.Ticker, bond.Ticker)
		}
	})

	t.Run("handles empty values", func(t *testing.T) {
		position := &PositionByFIFO{
			Ticker:             "SOME_VALUE",
			ClassCode:          "SOME_VALUE",
			Replaced:           true,
			CurrencyIfReplaced: "RUB",
		}

		bond := domain.BondIdentIdentifiers{
			Ticker:          "",
			ClassCode:       "",
			Replaced:        false,
			NominalCurrency: "",
		}

		position.ApplyBondMetadata(bond)

		// Проверяем, что поля стали пустыми
		if position.Ticker != "" || position.ClassCode != "" ||
			position.Replaced != false || position.CurrencyIfReplaced != "" {
			t.Errorf("fields not properly cleared: %+v", position)
		}
	})
}

func TestIsCurrentQuantityGreaterThanSellQuantity(t *testing.T) {
	t.Run("successfully reduces quantity and adjusts proportional fields", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:       100.0,
			TotalComission: 50.0,
			PaidTax:        200.0,
			BuyAccruedInt:  30.0,
		}
		sellQuantity := 30.0

		err := position.isCurrentQuantityGreaterThanSellQuantity(sellQuantity)
		// Проверяем, что ошибки нет
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Проверяем, что количество уменьшилось
		expectedQuantity := 70.0
		if position.Quantity != expectedQuantity {
			t.Errorf("Quantity = %f, want %f", position.Quantity, expectedQuantity)
		}

		// Проверяем пропорциональное уменьшение полей
		// Продано 30% (30/100), значит поля должны уменьшиться на 30%
		expectedComission := 35.0  // 50 * (1 - 0.3)
		expectedTax := 140.0       // 200 * (1 - 0.3)
		expectedAccruedInt := 21.0 // 30 * (1 - 0.3)

		if position.TotalComission != expectedComission {
			t.Errorf("TotalComission = %f, want %f", position.TotalComission, expectedComission)
		}
		if position.PaidTax != expectedTax {
			t.Errorf("PaidTax = %f, want %f", position.PaidTax, expectedTax)
		}
		if position.BuyAccruedInt != expectedAccruedInt {
			t.Errorf("BuyAccruedInt = %f, want %f", position.BuyAccruedInt, expectedAccruedInt)
		}
	})

	t.Run("sell all quantity", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:       100.0,
			TotalComission: 50.0,
			PaidTax:        200.0,
			BuyAccruedInt:  30.0,
		}
		sellQuantity := 100.0

		err := position.isCurrentQuantityGreaterThanSellQuantity(sellQuantity)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Все поля должны стать нулевыми (1 - 1 = 0)
		if position.Quantity != 0 {
			t.Errorf("Quantity = %f, want 0", position.Quantity)
		}
		if position.TotalComission != 0 {
			t.Errorf("TotalComission = %f, want 0", position.TotalComission)
		}
		if position.PaidTax != 0 {
			t.Errorf("PaidTax = %f, want 0", position.PaidTax)
		}
		if position.BuyAccruedInt != 0 {
			t.Errorf("BuyAccruedInt = %f, want 0", position.BuyAccruedInt)
		}
	})

	t.Run("sell quantity less than current", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:       100.0,
			TotalComission: 50.0,
			PaidTax:        200.0,
			BuyAccruedInt:  30.0,
		}
		sellQuantity := 50.0

		err := position.isCurrentQuantityGreaterThanSellQuantity(sellQuantity)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Продано 50%, поля должны уменьшиться вдвое
		if position.Quantity != 50.0 {
			t.Errorf("Quantity = %f, want 50.0", position.Quantity)
		}
		if position.TotalComission != 25.0 {
			t.Errorf("TotalComission = %f, want 25.0", position.TotalComission)
		}
		if position.PaidTax != 100.0 {
			t.Errorf("PaidTax = %f, want 100.0", position.PaidTax)
		}
		if position.BuyAccruedInt != 15.0 {
			t.Errorf("BuyAccruedInt = %f, want 15.0", position.BuyAccruedInt)
		}
	})

	t.Run("zero current quantity returns error", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity: 0,
		}
		sellQuantity := 10.0

		err := position.isCurrentQuantityGreaterThanSellQuantity(sellQuantity)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != ErrZeroQuantity {
			t.Errorf("expected ErrZeroQuantity, got %v", err)
		}
	})

	t.Run("sell quantity can be greater than current - function doesn't validate this", func(t *testing.T) {
		// Функция не проверяет, что sellQuantity <= currentQuantity
		// Это важно знать, если такое поведение ожидаемо или нет
		position := &PositionByFIFO{
			Quantity:       100.0,
			TotalComission: 50.0,
			PaidTax:        200.0,
			BuyAccruedInt:  30.0,
		}
		sellQuantity := 150.0

		err := position.isCurrentQuantityGreaterThanSellQuantity(sellQuantity)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Количество станет отрицательным
		if position.Quantity != -50.0 {
			t.Errorf("Quantity = %f, want -50.0", position.Quantity)
		}

		// Пропорция > 1, поля станут отрицательными
		expectedComission := -25.0  // 50 * (1 - 1.5)
		expectedTax := -100.0       // 200 * (1 - 1.5)
		expectedAccruedInt := -15.0 // 30 * (1 - 1.5)

		if position.TotalComission != expectedComission {
			t.Errorf("TotalComission = %f, want %f", position.TotalComission, expectedComission)
		}
		if position.PaidTax != expectedTax {
			t.Errorf("PaidTax = %f, want %f", position.PaidTax, expectedTax)
		}
		if position.BuyAccruedInt != expectedAccruedInt {
			t.Errorf("BuyAccruedInt = %f, want %f", position.BuyAccruedInt, expectedAccruedInt)
		}
	})

	t.Run("handles fractional quantities", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:       100.5,
			TotalComission: 50.25,
			PaidTax:        200.75,
			BuyAccruedInt:  30.6,
		}
		sellQuantity := 33.3

		err := position.isCurrentQuantityGreaterThanSellQuantity(sellQuantity)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Проверяем с допустимой погрешностью для float
		proportion := 33.3 / 100.5

		expectedQuantity := 100.5 - 33.3
		expectedComission := 50.25 * (1 - proportion)
		expectedTax := 200.75 * (1 - proportion)
		expectedAccruedInt := 30.6 * (1 - proportion)

		delta := 0.0001

		if abs(position.Quantity-expectedQuantity) > delta {
			t.Errorf("Quantity = %f, want %f", position.Quantity, expectedQuantity)
		}
		if abs(position.TotalComission-expectedComission) > delta {
			t.Errorf("TotalComission = %f, want %f", position.TotalComission, expectedComission)
		}
		if abs(position.PaidTax-expectedTax) > delta {
			t.Errorf("PaidTax = %f, want %f", position.PaidTax, expectedTax)
		}
		if abs(position.BuyAccruedInt-expectedAccruedInt) > delta {
			t.Errorf("BuyAccruedInt = %f, want %f", position.BuyAccruedInt, expectedAccruedInt)
		}
	})

	t.Run("other fields remain unchanged", func(t *testing.T) {
		position := &PositionByFIFO{
			Name:               "Test Bond",
			Replaced:           true,
			CurrencyIfReplaced: "USD",
			BuyDate:            time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Figi:               "BBG000000001",
			InstrumentType:     "Bond",
			Quantity:           100.0,
			TotalComission:     50.0,
			PaidTax:            200.0,
			BuyAccruedInt:      30.0,
		}
		sellQuantity := 30.0

		// Сохраняем значения полей, которые не должны измениться
		originalName := position.Name
		originalReplaced := position.Replaced
		originalCurrency := position.CurrencyIfReplaced
		originalBuyDate := position.BuyDate
		originalFigi := position.Figi
		originalInstrumentType := position.InstrumentType

		err := position.isCurrentQuantityGreaterThanSellQuantity(sellQuantity)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Проверяем, что эти поля не изменились
		if position.Name != originalName {
			t.Errorf("Name changed to %q", position.Name)
		}
		if position.Replaced != originalReplaced {
			t.Errorf("Replaced changed to %v", position.Replaced)
		}
		if position.CurrencyIfReplaced != originalCurrency {
			t.Errorf("CurrencyIfReplaced changed to %q", position.CurrencyIfReplaced)
		}
		if !position.BuyDate.Equal(originalBuyDate) {
			t.Errorf("BuyDate changed to %v", position.BuyDate)
		}
		if position.Figi != originalFigi {
			t.Errorf("Figi changed to %q", position.Figi)
		}
		if position.InstrumentType != originalInstrumentType {
			t.Errorf("InstrumentType changed to %q", position.InstrumentType)
		}
	})
}

func TestGetSecurityIncomeWithoutTax(t *testing.T) {
	t.Run("closed position with all components", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              10.0,
			BuyPrice:              100.0,
			SellPrice:             110.0,
			BuyAccruedInt:         5.0,
			SellAccruedInt:        8.0,
			TotalCoupon:           20.0,
			TotalDividend:         15.0,
			TotalComission:        2.5,
			PartialEarlyRepayment: 0.0,
		}

		// Расчет:
		// buySellDifference = (110-100)*10 + 8 - 5 = 100 + 3 = 103
		// cashFlow = 20 + 15 = 35
		// positionProfit = 103 + 35 + 2.5 + 0 = 140.5
		expected := 140.5

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("open position with potential sell commission", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              5.0,
			BuyPrice:              200.0,
			SellPrice:             210.0, // текущая рыночная цена
			BuyAccruedInt:         3.0,
			SellAccruedInt:        4.0,
			TotalCoupon:           10.0,
			TotalDividend:         0.0,
			TotalComission:        1.5,
			PartialEarlyRepayment: 0.0,
		}

		// Расчет:
		// buySellDifference = (210-200)*5 + 4 - 3 = 50 + 1 = 51
		// cashFlow = 10 + 0 = 10
		// positionProfit = 51 + 10 + 1.5 + 0 = 62.5
		expected := 62.5

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("position with partial early repayment", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              20.0,
			BuyPrice:              50.0,
			SellPrice:             55.0,
			BuyAccruedInt:         2.0,
			SellAccruedInt:        3.0,
			TotalCoupon:           0.0,
			TotalDividend:         0.0,
			TotalComission:        1.0,
			PartialEarlyRepayment: 30.0,
		}

		// Расчет:
		// buySellDifference = (55-50)*20 + 3 - 2 = 100 + 1 = 101
		// cashFlow = 0
		// positionProfit = 101 + 0 + 1.0 + 30.0 = 132.0
		expected := 132.0

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("position with negative difference (loss)", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              15.0,
			BuyPrice:              100.0,
			SellPrice:             90.0,
			BuyAccruedInt:         4.0,
			SellAccruedInt:        2.0,
			TotalCoupon:           12.0,
			TotalDividend:         0.0,
			TotalComission:        3.0,
			PartialEarlyRepayment: 0.0,
		}

		// Расчет:
		// buySellDifference = (90-100)*15 + 2 - 4 = -150 - 2 = -152
		// cashFlow = 12
		// positionProfit = -152 + 12 + 3.0 + 0 = -137
		expected := -137.0

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("zero quantity position", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              0.0,
			BuyPrice:              100.0,
			SellPrice:             110.0,
			BuyAccruedInt:         5.0,
			SellAccruedInt:        8.0,
			TotalCoupon:           20.0,
			TotalDividend:         15.0,
			TotalComission:        2.5,
			PartialEarlyRepayment: 10.0,
		}

		// Расчет с quantity = 0:
		// buySellDifference = (110-100)*0 + 8 - 5 = 3
		// cashFlow = 35
		// positionProfit = 3 + 35 + 2.5 + 10 = 50.5
		expected := 50.5

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("only coupon income, no price difference", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              8.0,
			BuyPrice:              100.0,
			SellPrice:             100.0,
			BuyAccruedInt:         5.0,
			SellAccruedInt:        5.0,
			TotalCoupon:           40.0,
			TotalDividend:         0.0,
			TotalComission:        2.0,
			PartialEarlyRepayment: 0.0,
		}

		// Расчет:
		// buySellDifference = (100-100)*8 + 5 - 5 = 0
		// cashFlow = 40
		// positionProfit = 0 + 40 + 2.0 + 0 = 42.0
		expected := 42.0

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("negative values in components", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              10.0,
			BuyPrice:              100.0,
			SellPrice:             90.0,
			BuyAccruedInt:         5.0,
			SellAccruedInt:        3.0,
			TotalCoupon:           -2.0, // отрицательный купон (теоретически)
			TotalDividend:         0.0,
			TotalComission:        -1.0, // отрицательная комиссия (возврат)
			PartialEarlyRepayment: 0.0,
		}

		// Расчет:
		// buySellDifference = (90-100)*10 + 3 - 5 = -100 - 2 = -102
		// cashFlow = -2
		// positionProfit = -102 + (-2) + (-1.0) + 0 = -105
		expected := -105.0

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("all fields zero", func(t *testing.T) {
		position := &PositionByFIFO{}

		result := position.GetProfitBeforeTax()

		if result != 0.0 {
			t.Errorf("expected 0, got %f", result)
		}
	})

	t.Run("fractional values precision", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              10.5,
			BuyPrice:              100.25,
			SellPrice:             110.75,
			BuyAccruedInt:         5.125,
			SellAccruedInt:        8.375,
			TotalCoupon:           20.625,
			TotalDividend:         15.875,
			TotalComission:        2.25,
			PartialEarlyRepayment: 0.0,
		}

		// Расчет с плавающей точкой
		buySellDiff := (110.75-100.25)*10.5 + 8.375 - 5.125
		cashFlow := 20.625 + 15.875
		expected := buySellDiff + cashFlow + 2.25

		result := position.GetProfitBeforeTax()

		// Проверяем с допустимой погрешностью
		delta := 0.0001
		if abs(result-expected) > delta {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("large numbers", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              1000000.0,
			BuyPrice:              1000.0,
			SellPrice:             1100.0,
			BuyAccruedInt:         50000.0,
			SellAccruedInt:        80000.0,
			TotalCoupon:           200000.0,
			TotalDividend:         150000.0,
			TotalComission:        25000.0,
			PartialEarlyRepayment: 100000.0,
		}

		// Расчет:
		// buySellDifference = (1100-1000)*1e6 + 80000 - 50000 = 100e6 + 30000 = 100030000
		// cashFlow = 200000 + 150000 = 350000
		// positionProfit = 100030000 + 350000 + 25000 + 100000 = 100505000
		expected := 100505000.0

		result := position.GetProfitBeforeTax()

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("verify formula components", func(t *testing.T) {
		position := &PositionByFIFO{
			Quantity:              10.0,
			BuyPrice:              100.0,
			SellPrice:             110.0,
			BuyAccruedInt:         5.0,
			SellAccruedInt:        8.0,
			TotalCoupon:           20.0,
			TotalDividend:         15.0,
			TotalComission:        2.5,
			PartialEarlyRepayment: 0.0,
		}

		result := position.GetProfitBeforeTax()

		// Разбиваем на компоненты для проверки
		buySellDiff := (position.SellPrice-position.BuyPrice)*position.Quantity +
			position.SellAccruedInt - position.BuyAccruedInt
		cashFlow := position.TotalCoupon + position.TotalDividend
		manualSum := buySellDiff + cashFlow + position.TotalComission + position.PartialEarlyRepayment

		if result != manualSum {
			t.Errorf("result %f doesn't match manual sum %f", result, manualSum)
		}
	})
}

func TestIsHoldingPeriodMoreThanThreeYears(t *testing.T) {
	t.Run("продажа через 2 года - должно быть false", func(t *testing.T) {
		buyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != false {
			t.Errorf("expected false for 2 years holding, got true")
		}
	})

	t.Run("продажа через 3 года ровно - должно быть true", func(t *testing.T) {
		buyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != true {
			t.Errorf("expected true for exactly 3 years holding, got false")
		}
	})

	t.Run("продажа через 3 года и 1 день - должно быть true", func(t *testing.T) {
		buyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != true {
			t.Errorf("expected true for 3 years + 1 day holding, got false")
		}
	})

	t.Run("продажа через 10 лет - должно быть true", func(t *testing.T) {
		buyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != true {
			t.Errorf("expected true for 10 years holding, got false")
		}
	})

	t.Run("продажа в тот же день - должно быть false", func(t *testing.T) {
		buyDate := time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC)
		sellDate := time.Date(2023, 6, 15, 15, 45, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != false {
			t.Errorf("expected false for same day sale, got true")
		}
	})

	t.Run("продажа через 2 года 364 дня - должно быть false", func(t *testing.T) {
		buyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != false {
			t.Errorf("expected false for 2 years 364 days holding, got true")
		}
	})

	t.Run("продажа через 3 года минус 1 день - должно быть false", func(t *testing.T) {
		buyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != false {
			t.Errorf("expected false for 3 years minus 1 day, got true")
		}
	})

	t.Run("продажа через 3 года с учетом високосного года", func(t *testing.T) {
		buyDate := time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC) // високосный год
		sellDate := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC) // после 3 лет

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		// AddDate(3,0,0) к 2020-02-29 даст 2023-03-01
		// sellDate 2023-03-01 >= 2023-03-01 -> true
		if result != true {
			t.Errorf("expected true for leap year case, got false")
		}
	})

	t.Run("продажа раньше покупки - должно быть false", func(t *testing.T) {
		buyDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		// Технически, если sellDate раньше buyDate, то threeYearAfterBuyDate.After(sellDate) = true
		// значит функция вернет false
		if result != false {
			t.Errorf("expected false for sell before buy, got true")
		}
	})

	t.Run("граница месяца: купили 31 января", func(t *testing.T) {
		buyDate := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

		tests := []struct {
			name     string
			sellDate time.Time
			expected bool
		}{
			{
				name:     "продажа 30 января 2023",
				sellDate: time.Date(2023, 1, 30, 0, 0, 0, 0, time.UTC),
				expected: false, // еще не 3 года
			},
			{
				name:     "продажа 31 января 2023",
				sellDate: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				expected: true, // ровно 3 года
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := isHoldingPeriodMoreThanThreeYears(buyDate, tt.sellDate)
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			})
		}
	})

	t.Run("разные часовые пояса", func(t *testing.T) {
		loc := time.FixedZone("UTC+3", 3*60*60)
		buyDate := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)
		sellDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		// Должно работать независимо от timezone
		if result != true {
			t.Errorf("expected true across timezones, got false")
		}
	})

	t.Run("очень старые даты", func(t *testing.T) {
		buyDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		sellDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

		result := isHoldingPeriodMoreThanThreeYears(buyDate, sellDate)

		if result != true {
			t.Errorf("expected true for century-spanning holding, got false")
		}
	})
}

func TestGetTotalTaxFromPosition(t *testing.T) {
	t.Run("прибыль положительная, срок менее 3 лет - налог 13%", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			SellDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		profit := 1000.0
		expected := 130.0 // 1000 * 0.13

		result := position.GetTotalTaxFromPosition(profit)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("прибыль положительная, срок ровно 3 года - налог 0", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			SellDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		profit := 1000.0
		expected := 0.0

		result := position.GetTotalTaxFromPosition(profit)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("прибыль положительная, срок более 3 лет - налог 0", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			SellDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		profit := 1000.0
		expected := 0.0

		result := position.GetTotalTaxFromPosition(profit)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("прибыль отрицательная (убыток) - налог 0 независимо от срока", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			SellDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		profit := -500.0
		expected := 0.0

		result := position.GetTotalTaxFromPosition(profit)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("прибыль ноль - налог 0", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			SellDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		profit := 0.0
		expected := 0.0

		result := position.GetTotalTaxFromPosition(profit)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("открытая позиция (нет SellDate) - функция использует zero value", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			// SellDate не указан - будет zero value (0001-01-01)
		}
		profit := 1000.0

		// С zero value SellDate срок владения будет отрицательным
		// isHoldingPeriodMoreThanThreeYears вернет false
		// Значит налог будет рассчитан
		expected := 130.0

		result := position.GetTotalTaxFromPosition(profit)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("некорректные даты (продажа раньше покупки)", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			SellDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		profit := 1000.0

		// С отрицательным сроком isHoldingPeriodMoreThanThreeYears вернет false
		// Значит налог будет рассчитан
		expected := 130.0

		result := position.GetTotalTaxFromPosition(profit)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("проверка с дробной прибылью", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			SellDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		profit := 1234.56
		expected := 160.4928 // 1234.56 * 0.13

		result := position.GetTotalTaxFromPosition(profit)

		// Проверяем с допустимой погрешностью для float
		delta := 0.0001
		if abs(result-expected) > delta {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})
}

func TestGetAnnualizedReturnInPercentage(t *testing.T) {
	t.Run("ровно 1 год с прибылью", func(t *testing.T) {
		buyDate := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)
		sellDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := 870.0

		expected := 8.7

		result, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		require.NoError(t, err)
		if math.Abs(result-expected) > 0.02 {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("6 месяцев с прибылью", func(t *testing.T) {
		buyDate := time.Date(2023, 7, 15, 12, 0, 0, 0, time.UTC)
		sellDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := 435.0

		// Ожидаем ~8.9% (4.35% за полгода в годовом выражении)
		expected := 8.9

		result, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if math.Abs(result-expected) > 0.1 {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("2 года с прибылью", func(t *testing.T) {
		buyDate := time.Date(2022, 1, 15, 12, 0, 0, 0, time.UTC)
		sellDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := 1740.0

		// 17.4% за 2 года = ~8.36% годовых
		expected := 8.36

		result, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if math.Abs(result-expected) > 0.1 {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("убыток", func(t *testing.T) {
		buyDate := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)
		sellDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := -500.0

		result, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result >= 0 {
			t.Errorf("expected negative return, got %f%%", result)
		}
	})

	t.Run("нулевая прибыль", func(t *testing.T) {
		buyDate := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)
		sellDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := 0.0

		result, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != 0 {
			t.Errorf("expected 0%%, got %f%%", result)
		}
	})

	t.Run("продажа в тот же день", func(t *testing.T) {
		buyDate := time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)
		sellDate := time.Date(2024, 1, 15, 16, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := 100.0

		result, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result < 1000 {
			t.Errorf("expected very high annualized return (>1000%%), got %f%%", result)
		}
	})

	t.Run("деление на ноль (buyPrice = 0)", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 0.0,
			Quantity: 10.0,
			BuyDate:  time.Now(),
		}

		_, err := position.GetAnnualizedReturnInPercentage(1000.0, time.Now())

		if err == nil {
			t.Error("expected error for zero buy price, got nil")
		}
		if err != ErrZeroDivision {
			t.Errorf("expected ErrZeroDivision, got %v", err)
		}
	})

	t.Run("деление на ноль (quantity = 0)", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 0.0,
			BuyDate:  time.Now(),
		}

		_, err := position.GetAnnualizedReturnInPercentage(1000.0, time.Now())

		if err == nil {
			t.Error("expected error for zero quantity, got nil")
		}
	})

	t.Run("продажа раньше покупки", func(t *testing.T) {
		buyDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
		sellDate := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := 870.0

		_, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrInvalidDate.Error())
	})

	t.Run("проверка округления до 2 знаков", func(t *testing.T) {
		buyDate := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)
		sellDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
			BuyDate:  buyDate,
		}
		netProfit := 873.21

		result, err := position.GetAnnualizedReturnInPercentage(netProfit, sellDate)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		rounded := math.Round(result*100) / 100
		if math.Abs(result-rounded) > 0.0001 {
			t.Errorf("result not rounded to 2 decimals: %f", result)
		}
	})
}

func TestGetProfit(t *testing.T) {
	t.Run("положительная прибыль", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}
		profit := 870.0

		// 870 / (1000 * 10) * 100 = 870 / 10000 * 100 = 8.7%
		expected := 8.7

		result, err := position.GetProfit(profit)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("отрицательная прибыль (убыток)", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}
		profit := -500.0

		// -500 / 10000 * 100 = -5%
		expected := -5.0

		result, err := position.GetProfit(profit)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("нулевая прибыль", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}
		profit := 0.0

		expected := 0.0

		result, err := position.GetProfit(profit)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("дробная прибыль с округлением", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}
		profit := 873.21

		// 873.21 / 10000 * 100 = 8.7321% -> округление до 8.73
		expected := 8.73

		result, err := position.GetProfit(profit)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("очень маленькая прибыль", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}
		profit := 0.01

		// 0.01 / 10000 * 100 = 0.0001% -> округление до 0.00
		expected := 0.0

		result, err := position.GetProfit(profit)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("очень большая прибыль", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}
		profit := 50000.0

		// 50000 / 10000 * 100 = 500%
		expected := 500.0

		result, err := position.GetProfit(profit)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("expected %f%%, got %f%%", expected, result)
		}
	})

	t.Run("деление на ноль (buyPrice = 0)", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 0.0,
			Quantity: 10.0,
		}

		_, err := position.GetProfit(1000.0)

		if err == nil {
			t.Error("expected error for zero buy price, got nil")
		}
		if err != ErrZeroDivision {
			t.Errorf("expected ErrZeroDivision, got %v", err)
		}
	})

	t.Run("деление на ноль (quantity = 0)", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 0.0,
		}

		_, err := position.GetProfit(1000.0)

		if err == nil {
			t.Error("expected error for zero quantity, got nil")
		}
	})

	t.Run("разные buyPrice и quantity с одинаковым totalInvest", func(t *testing.T) {
		position1 := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}

		position2 := &PositionByFIFO{
			BuyPrice: 500.0,
			Quantity: 20.0,
		}

		profit := 870.0

		result1, _ := position1.GetProfit(profit)
		result2, _ := position2.GetProfit(profit)

		if result1 != result2 {
			t.Errorf("expected equal results: %f vs %f", result1, result2)
		}
	})

	t.Run("проверка округления граничных значений", func(t *testing.T) {
		position := &PositionByFIFO{
			BuyPrice: 1000.0,
			Quantity: 10.0,
		}

		testCases := []struct {
			profit   float64
			expected float64
		}{
			{874.99, 8.75}, // 8.7499% -> 8.75
			{875.00, 8.75}, // 8.75% -> 8.75
			{875.01, 8.75}, // 8.7501% -> 8.75 (если RoundFloat банковское)
		}

		for _, tc := range testCases {
			result, err := position.GetProfit(tc.profit)
			if err != nil {
				t.Errorf("unexpected error for profit %f: %v", tc.profit, err)
			}
			if result != tc.expected {
				t.Errorf("profit %f: expected %f%%, got %f%%", tc.profit, tc.expected, result)
			}
		}
	})
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
