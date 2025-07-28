package visualization

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"main.go/service/service_models"
)

func Vizualize(reports []service_models.BondReport, filename string) error {
	const (
		width        = 1200
		heightPerRow = 40
		margin       = 20
		headerHeight = 60
		rowHeight    = 30
		colCount     = 8 // Количество отображаемых колонок
	)

	// Рассчитываем высоту изображения
	height := headerHeight + len(reports)*rowHeight + margin*2
	if height < 300 {
		height = 300
	}

	// Создаем контекст изображения
	dc := gg.NewContext(width, height)

	// Загружаем шрифт
	font, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})
	if err != nil {
		return err
	}
	dc.SetFontFace(face)

	// Заливаем фон
	dc.SetColor(color.White)
	dc.Clear()

	// Рисуем заголовок
	dc.SetColor(color.Black)
	dc.DrawStringAnchored("Отчет по облигациям", width/2, margin, 0.5, 0.5)
	dc.SetFontFace(face)

	// Определяем колонки для отображения
	columns := []struct {
		Title string
		Width float64
	}{
		{"Тикер", 80},
		{"Название", 200},
		{"Дата погаш.", 100},
		{"Доходн. к погаш.", 120},
		{"Тек. цена", 100},
		{"Номинал", 100},
		{"Прибыль", 100},
		{"Доходн. год.", 120},
	}

	// Рисуем заголовки таблицы
	y := float64(headerHeight)
	for i, col := range columns {
		x := margin + float64(i)*col.Width
		dc.SetColor(color.RGBA{200, 200, 200, 255})
		dc.DrawRectangle(x, y, col.Width, rowHeight)
		dc.Fill()
		dc.SetColor(color.Black)
		dc.DrawStringAnchored(col.Title, x+col.Width/2, y+rowHeight/2, 0.5, 0.5)
	}

	// Рисуем данные
	for _, report := range reports {
		y += rowHeight
		dc.SetColor(color.Black)

		// Форматируем данные
		values := []string{
			report.Ticker,
			report.Name,
			report.MaturityDate,
			formatPercent(report.YieldToMaturity),
			formatCurrency(report.CurrentPrice),
			formatCurrency(report.Nominal),
			formatCurrency(report.Profit),
			formatPercent(report.AnnualizedReturn),
		}

		for i, value := range values {
			x := margin + float64(i)*columns[i].Width
			dc.DrawStringAnchored(value, x+columns[i].Width/2, y+rowHeight/2, 0.5, 0.5)
		}
	}

	// Сохраняем изображение
	return dc.SavePNG(filename)
}

func formatCurrency(value float64) string {
	return fmt.Sprintf("%.2f ₽", value)
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value*100)
}
