package report

import (
	"bonds-report-service/internal/domain"
	"testing"
	"time"
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

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
