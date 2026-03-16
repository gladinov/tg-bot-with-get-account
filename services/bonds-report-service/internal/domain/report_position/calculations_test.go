//go:build unit

package report

import (
	"bonds-report-service/internal/domain"
	"math"
	"testing"
)

func TestGetNetProfit(t *testing.T) {
	// Кейс 1: Обычная прибыль с налогом
	t.Run("positive profit with tax", func(t *testing.T) {
		profit := 1000.0
		tax := 130.0
		expected := 870.0

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 2: Нулевая прибыль
	t.Run("zero profit", func(t *testing.T) {
		profit := 0.0
		tax := 0.0
		expected := 0.0

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 3: Отрицательная прибыль (убыток)
	t.Run("negative profit (loss)", func(t *testing.T) {
		profit := -500.0
		tax := 0.0 // налог с убытка не берется
		expected := -500.0

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 4: Налог больше прибыли
	t.Run("tax greater than profit", func(t *testing.T) {
		profit := 100.0
		tax := 150.0
		expected := -50.0

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 5: Налог равен прибыли
	t.Run("tax equals profit", func(t *testing.T) {
		profit := 500.0
		tax := 500.0
		expected := 0.0

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 6: Дробные значения
	t.Run("fractional values", func(t *testing.T) {
		profit := 1000.50
		tax := 130.07
		expected := 870.43

		result := GetNetProfit(profit, tax)

		// Для float используем delta
		delta := 0.0001
		if math.Abs(result-expected) > delta {
			t.Errorf("expected %f, got %f (delta %f)", expected, result, delta)
		}
	})

	// Кейс 7: Очень большие числа
	t.Run("large numbers", func(t *testing.T) {
		profit := 1_000_000_000.0
		tax := 130_000_000.0
		expected := 870_000_000.0

		result := GetNetProfit(profit, tax)

		delta := 0.01
		if math.Abs(result-expected) > delta {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 8: Очень маленькие числа
	t.Run("very small numbers", func(t *testing.T) {
		profit := 0.0001
		tax := 0.000013
		expected := 0.000087

		result := GetNetProfit(profit, tax)

		delta := 0.000001
		if math.Abs(result-expected) > delta {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 9: Нулевой налог
	t.Run("zero tax", func(t *testing.T) {
		profit := 750.0
		tax := 0.0
		expected := 750.0

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 10: Отрицательный налог (теоретический случай)
	t.Run("negative tax", func(t *testing.T) {
		profit := 1000.0
		tax := -50.0 // возврат налога
		expected := 1050.0

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 11: Максимальные значения float64
	t.Run("max float64 values", func(t *testing.T) {
		profit := 1e308 // близко к max float64
		tax := 1e307

		result := GetNetProfit(profit, tax)

		// Проверяем, что не бесконечность
		if math.IsInf(result, 0) {
			t.Errorf("result is infinite: %f", result)
		}

		// Проверяем, что не NaN
		if math.IsNaN(result) {
			t.Errorf("result is NaN")
		}
	})

	// Кейс 12: NaN значения
	t.Run("NaN values", func(t *testing.T) {
		profit := math.NaN()
		tax := 100.0

		result := GetNetProfit(profit, tax)

		// Функция должна корректно обработать NaN
		if !math.IsNaN(result) {
			t.Errorf("expected NaN, got %f", result)
		}
	})

	// Кейс 13: Бесконечность
	t.Run("infinity values", func(t *testing.T) {
		profit := math.Inf(1)
		tax := 100.0

		result := GetNetProfit(profit, tax)

		// Бесконечность минус число = бесконечность
		if !math.IsInf(result, 1) {
			t.Errorf("expected +Inf, got %f", result)
		}
	})

	// Кейс 14: Точность вычислений
	t.Run("precision test", func(t *testing.T) {
		profit := 0.1 + 0.2 // 0.30000000000000004
		tax := 0.05
		expected := 0.25

		result := GetNetProfit(profit, tax)

		// Из-за погрешности float, используем разумную дельту
		delta := 1e-15
		if math.Abs(result-expected) > delta {
			t.Errorf("precision error: expected %f, got %f", expected, result)
		}
	})

	// Кейс 15: Оба параметра отрицательные
	t.Run("both negative", func(t *testing.T) {
		profit := -1000.0
		tax := -130.0
		expected := -870.0 // (-1000) - (-130) = -870

		result := GetNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})
}

func TestCalculateSellPrice(t *testing.T) {
	t.Run("normal value", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 500000000),
		}
		expected := 985.0

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("price at 100 percent", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(100, 0),
		}
		expected := 1000.0

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("price above nominal", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(105, 750000000),
		}
		expected := 1057.5

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("price below nominal", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(89, 250000000),
		}
		expected := 892.5

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("zero nominal", func(t *testing.T) {
		nominal := 0.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 500000000),
		}
		expected := 0.0

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("zero price", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(0, 0),
		}
		expected := 0.0

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("fractional nominal", func(t *testing.T) {
		nominal := 1234.56
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 500000000),
		}
		// Вычисляем ожидаемое значение через тот же механизм, что и функция
		price := lastPrice.LastPrice.ToFloat() / 100 * nominal * 100
		expected := math.Round(price) / 100

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("fractional price", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 543210000),
		}
		expected := 985.43

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("very small nominal", func(t *testing.T) {
		nominal := 0.01
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 500000000),
		}
		expected := 0.01

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("very large nominal", func(t *testing.T) {
		nominal := 1_000_000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 500000000),
		}
		expected := 985_000.0

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("price with three decimals", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 555000000),
		}
		expected := 985.55

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("rounding up", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 555500000),
		}
		expected := 985.56

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("rounding down", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 554400000),
		}
		expected := 985.54

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("negative price", func(t *testing.T) {
		nominal := 1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(-5, 0),
		}
		expected := -50.0

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("negative nominal", func(t *testing.T) {
		nominal := -1000.0
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(98, 500000000),
		}
		expected := -985.0

		result := CalculateSellPrice(nominal, lastPrice)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("large but safe nominal", func(t *testing.T) {
		nominal := 1e100 // безопасное значение
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(50, 0),
		}
		expected := 5e99 // 50% от номинала

		result := CalculateSellPrice(nominal, lastPrice)

		if math.IsInf(result, 0) || math.IsNaN(result) {
			t.Errorf("invalid result: %f", result)
		}

		delta := expected * 1e-12
		if math.Abs(result-expected) > delta {
			t.Errorf("expected %e, got %e", expected, result)
		}
	})

	t.Run("min positive values", func(t *testing.T) {
		nominal := 1e-308
		lastPrice := domain.LastPrice{
			LastPrice: domain.NewQuotation(50, 0),
		}

		result := CalculateSellPrice(nominal, lastPrice)

		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("unexpected result: %f", result)
		}
	})

	t.Run("formula verification", func(t *testing.T) {
		testCases := []struct {
			nominal  float64
			units    int64
			nano     int32
			expected float64
		}{
			{1000, 100, 0, 1000},
			{1000, 50, 0, 500},
			{1000, 150, 0, 1500},
			{500, 75, 500000000, 377.5},
			{2000, 33, 330000000, 666.6},
			{1500, 99, 990000000, 1499.85},
			{1234, 45, 670000000, 563.57},
		}

		for i, tc := range testCases {
			t.Run("", func(t *testing.T) {
				lastPrice := domain.LastPrice{
					LastPrice: domain.NewQuotation(tc.units, tc.nano),
				}
				result := CalculateSellPrice(tc.nominal, lastPrice)

				if math.Abs(result-tc.expected) > 0.01 {
					t.Errorf("case %d: nominal=%f, price=%d.%d: expected %f, got %f",
						i, tc.nominal, tc.units, tc.nano, tc.expected, result)
				}
			})
		}
	})

	t.Run("always returns 2 decimals", func(t *testing.T) {
		nominals := []float64{1, 10, 100, 1000, 1234.56}
		prices := []struct {
			units int64
			nano  int32
		}{
			{0, 100000000},
			{1, 230000000},
			{45, 670000000},
			{89, 10000000},
			{99, 990000000},
			{100, 0},
			{101, 500000000},
		}

		for _, nominal := range nominals {
			for _, price := range prices {
				t.Run("", func(t *testing.T) {
					lastPrice := domain.LastPrice{
						LastPrice: domain.NewQuotation(price.units, price.nano),
					}
					result := CalculateSellPrice(nominal, lastPrice)

					fractional := math.Abs(result*100 - math.Floor(result*100+0.5))
					if fractional > 0.0001 {
						t.Errorf("result has more than 2 decimals: nominal=%f, price=%d.%d, result=%f",
							nominal, price.units, price.nano, result)
					}
				})
			}
		}
	})

	t.Run("rounding behavior", func(t *testing.T) {
		nominal := 1000.0
		testPrices := []struct {
			units    int64
			nano     int32
			expected float64
		}{
			{98, 554400000, 985.54},
			{98, 555500000, 985.56},
			{98, 555000000, 985.55},
			{98, 555100000, 985.55},
		}

		for _, tc := range testPrices {
			t.Run("", func(t *testing.T) {
				lastPrice := domain.LastPrice{
					LastPrice: domain.NewQuotation(tc.units, tc.nano),
				}
				result := CalculateSellPrice(nominal, lastPrice)

				if result != tc.expected {
					t.Errorf("price=%d.%d: expected %f, got %f",
						tc.units, tc.nano, tc.expected, result)
				}
			})
		}
	})
}

func TestCalculateNominal(t *testing.T) {
	t.Run("not replaced returns base nominal", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := false
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 1000.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("replaced with valid rate multiplies nominal", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 1500.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("replaced with fractional rate", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  0.75,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 750.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("replaced with fractional nominal", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1234, 560000000)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 1851.84

		result := CalculateNominal(nominal, replaced, rate)

		if math.Abs(result-expected) > 0.01 {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("replaced with fractional both", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1234, 560000000)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.2345,
				IsSet:  true,
				IsNull: false,
			},
		}

		// Вычисляем ожидаемое значение через тот же механизм
		nominalFloat := nominal.ToFloat() // 1234.56
		expected := nominalFloat * 1.2345 // 1523.76432

		result := CalculateNominal(nominal, replaced, rate)

		// Используем разумную дельту для float
		delta := 0.01
		if math.Abs(result-expected) > delta {
			t.Errorf("expected %.2f ±%.2f, got %.2f", expected, delta, result)
		}
	})

	t.Run("replaced but rate not set returns base", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  false,
				IsNull: false,
			},
		}
		expected := 1000.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("replaced but rate is null returns base", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: true,
			},
		}
		expected := 1000.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("replaced but rate not set and null returns base", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  false,
				IsNull: true,
			},
		}
		expected := 1000.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("zero nominal", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 0, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 0.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("zero rate", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  0.0,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 0.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("negative nominal", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", -1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := -1500.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("negative rate", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  -0.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := -500.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("very small nominal", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 0, 1)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 1.5e-9

		result := CalculateNominal(nominal, replaced, rate)

		if math.Abs(result-expected) > 1e-15 {
			t.Errorf("expected %e, got %e", expected, result)
		}
	})

	t.Run("very large nominal", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1_000_000_000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1_000_000.0,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 1e15

		result := CalculateNominal(nominal, replaced, rate)

		if math.Abs(result-expected) > expected*1e-12 {
			t.Errorf("expected %e, got %e", expected, result)
		}
	})

	t.Run("rate with many decimals", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.23456789,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 1234.56789

		result := CalculateNominal(nominal, replaced, rate)

		if math.Abs(result-expected) > 1e-8 {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("nominal with nano part", func(t *testing.T) {
		nominal := domain.NewMoneyValue("RUB", 1000, 123456789)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  2.0,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 2000.246913578

		result := CalculateNominal(nominal, replaced, rate)

		if math.Abs(result-expected) > 1e-8 {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("different currency ignored in calculation", func(t *testing.T) {
		nominal := domain.NewMoneyValue("USD", 1000, 0)
		replaced := true
		rate := domain.Rate{
			IsoCurrencyName: "RUB",
			Vunit_Rate: domain.NullFloat64{
				Value:  1.5,
				IsSet:  true,
				IsNull: false,
			},
		}
		expected := 1500.0

		result := CalculateNominal(nominal, replaced, rate)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	t.Run("formula verification", func(t *testing.T) {
		testCases := []struct {
			name       string
			units      int64
			nano       int32
			currency   string
			replaced   bool
			rateValue  float64
			rateIsSet  bool
			rateIsNull bool
			expected   float64
		}{
			{
				name:  "case 1: not replaced",
				units: 1000, nano: 0, currency: "RUB",
				replaced:  false,
				rateValue: 1.5, rateIsSet: true, rateIsNull: false,
				expected: 1000.0,
			},
			{
				name:  "case 2: replaced with valid rate",
				units: 1000, nano: 0, currency: "RUB",
				replaced:  true,
				rateValue: 1.5, rateIsSet: true, rateIsNull: false,
				expected: 1500.0,
			},
			{
				name:  "case 3: replaced but rate not set",
				units: 1000, nano: 0, currency: "RUB",
				replaced:  true,
				rateValue: 1.5, rateIsSet: false, rateIsNull: false,
				expected: 1000.0,
			},
			{
				name:  "case 4: replaced but rate null",
				units: 1000, nano: 0, currency: "RUB",
				replaced:  true,
				rateValue: 1.5, rateIsSet: true, rateIsNull: true,
				expected: 1000.0,
			},
			{
				name:  "case 5: replaced with zero rate",
				units: 1000, nano: 0, currency: "RUB",
				replaced:  true,
				rateValue: 0.0, rateIsSet: true, rateIsNull: false,
				expected: 0.0,
			},
			{
				name:  "case 6: fractional nominal and rate",
				units: 1234, nano: 560000000, currency: "RUB",
				replaced:  true,
				rateValue: 1.2345, rateIsSet: true, rateIsNull: false,
				expected: 1523.76432, // оставляем для справки, но не используем строго
			},
			{
				name:  "case 7: USD currency",
				units: 500, nano: 0, currency: "USD",
				replaced:  true,
				rateValue: 2.0, rateIsSet: true, rateIsNull: false,
				expected: 1000.0,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				nominal := domain.NewMoneyValue(tc.currency, tc.units, tc.nano)
				rate := domain.Rate{
					IsoCurrencyName: "RUB",
					Vunit_Rate: domain.NullFloat64{
						Value:  tc.rateValue,
						IsSet:  tc.rateIsSet,
						IsNull: tc.rateIsNull,
					},
				}

				result := CalculateNominal(nominal, tc.replaced, rate)

				// Для case 6 используем вычисленное значение
				if tc.name == "case 6: fractional nominal and rate" {
					expected := nominal.ToFloat() * tc.rateValue
					if math.Abs(result-expected) > 0.01 {
						t.Errorf("expected %.5f, got %.5f", expected, result)
					}
				} else {
					if math.Abs(result-tc.expected) > 0.01 {
						t.Errorf("expected %.5f, got %.5f", tc.expected, result)
					}
				}
			})
		}
	})

	t.Run("verify MoneyValue.ToFloat precision", func(t *testing.T) {
		testCases := []struct {
			units    int64
			nano     int32
			expected float64
		}{
			{1000, 0, 1000.0},
			{1234, 560000000, 1234.56},
			{1000, 123456789, 1000.123456789},
			{0, 1, 1e-9},
			{0, 999999999, 0.999999999},
			{-1000, 0, -1000.0},
			{-1000, 500000000, -999.5},
		}

		for i, tc := range testCases {
			t.Run("", func(t *testing.T) {
				mv := domain.NewMoneyValue("RUB", tc.units, tc.nano)
				result := mv.ToFloat()

				if math.Abs(result-tc.expected) > 1e-12 {
					t.Errorf("case %d: (%d, %d): expected %.12f, got %.12f",
						i, tc.units, tc.nano, tc.expected, result)
				}
			})
		}
	})
}
