//go:build unit

package report

import (
	"math"
	"testing"
)

func TestGetNetProfit(t *testing.T) {
	// Кейс 1: Обычная прибыль с налогом
	t.Run("positive profit with tax", func(t *testing.T) {
		profit := 1000.0
		tax := 130.0
		expected := 870.0

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 2: Нулевая прибыль
	t.Run("zero profit", func(t *testing.T) {
		profit := 0.0
		tax := 0.0
		expected := 0.0

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 3: Отрицательная прибыль (убыток)
	t.Run("negative profit (loss)", func(t *testing.T) {
		profit := -500.0
		tax := 0.0 // налог с убытка не берется
		expected := -500.0

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 4: Налог больше прибыли
	t.Run("tax greater than profit", func(t *testing.T) {
		profit := 100.0
		tax := 150.0
		expected := -50.0

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 5: Налог равен прибыли
	t.Run("tax equals profit", func(t *testing.T) {
		profit := 500.0
		tax := 500.0
		expected := 0.0

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 6: Дробные значения
	t.Run("fractional values", func(t *testing.T) {
		profit := 1000.50
		tax := 130.07
		expected := 870.43

		result := getNetProfit(profit, tax)

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

		result := getNetProfit(profit, tax)

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

		result := getNetProfit(profit, tax)

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

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 10: Отрицательный налог (теоретический случай)
	t.Run("negative tax", func(t *testing.T) {
		profit := 1000.0
		tax := -50.0 // возврат налога
		expected := 1050.0

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})

	// Кейс 11: Максимальные значения float64
	t.Run("max float64 values", func(t *testing.T) {
		profit := 1e308 // близко к max float64
		tax := 1e307

		result := getNetProfit(profit, tax)

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

		result := getNetProfit(profit, tax)

		// Функция должна корректно обработать NaN
		if !math.IsNaN(result) {
			t.Errorf("expected NaN, got %f", result)
		}
	})

	// Кейс 13: Бесконечность
	t.Run("infinity values", func(t *testing.T) {
		profit := math.Inf(1)
		tax := 100.0

		result := getNetProfit(profit, tax)

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

		result := getNetProfit(profit, tax)

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

		result := getNetProfit(profit, tax)

		if result != expected {
			t.Errorf("expected %f, got %f", expected, result)
		}
	})
}
