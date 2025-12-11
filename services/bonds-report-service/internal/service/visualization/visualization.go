package visualization

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strings"
	"time"

	"bonds-report-service/internal/service/service_models"
	"bonds-report-service/lib/e"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	layout = "2006-01-02"
)

func Vizualize(reports []service_models.GeneralBondReportPosition, filename string, typeOfBonds string) error {

	const (
		width        = 1200 // Увеличим ширину для лучшего отображения
		heightPerRow = 40
		margin       = 20.0
		headerHeight = 60
		rowHeight    = 30
		colCount     = 15 // Количество отображаемых колонок
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
	var name string
	switch typeOfBonds {
	case service_models.ReplacedBonds:
		name = "Отчет по замещающим облигациям"
	case service_models.RubBonds:
		name = "Отчет по рублевым облигациям"
	case service_models.EuroBonds:
		name = "Отчет по валютным облигациям"
	}
	dc.DrawStringAnchored(name, width/2, margin, 0.5, 0.5)
	dc.SetFontFace(face)

	// Определяем колонки для отображения с более подходящими ширинами
	columns := []struct {
		Title string
		Width float64
		Flex  float64
	}{
		{"Тикер", 120, 0},              // 1
		{"Название", 180, 2},           // 2
		{"Валюта", 60, 0},              // 3
		{"Кол-во", 70, 0},              // 4
		{"% от портфеля", 90, 0},       // 5
		{"Дата погашения", 110, 0},     // 6
		{"Дюрация", 80, 0},             // 7
		{"Дата покупки", 90, 0},        // 8
		{"Ср.цена", 80, 0},             // 9
		{"Дох-ть при покупке", 120, 0}, // 10
		{"Дох-ть тек.", 90, 0},         // 11 (сокращенный заголовок)
		{"Тек.цена", 80, 0},            // 12 (сокращенный заголовок)
		{"Номинал", 80, 0},             // 13
		{"Доход", 80, 0},               // 14
		{"Дох-ть в %", 80, 0},          // 15
	}

	// Проверяем, что сумма минимальных ширин не превышает доступную ширину
	totalMinWidth := 0.0
	for _, col := range columns {
		totalMinWidth += col.Width
	}

	availableWidth := float64(width) - 2*margin
	if totalMinWidth > availableWidth {
		// Если минимальные ширины не помещаются, масштабируем их
		scaleFactor := availableWidth / totalMinWidth
		for i := range columns {
			columns[i].Width *= scaleFactor
		}
	} else {
		// Распределяем оставшееся пространство пропорционально Flex
		remainingWidth := availableWidth - totalMinWidth
		totalFlex := 0.0
		for _, col := range columns {
			totalFlex += col.Flex
		}

		if totalFlex > 0 {
			flexUnit := remainingWidth / totalFlex
			for i := range columns {
				if columns[i].Flex > 0 {
					columns[i].Width += flexUnit * columns[i].Flex
				}
			}
		}
	}

	// Рисование таблицы
	y := float64(headerHeight)
	currentX := margin

	// Заголовки
	for _, col := range columns {
		dc.SetColor(color.RGBA{200, 200, 200, 255})
		dc.DrawRectangle(currentX, y, col.Width, rowHeight)
		dc.Fill()

		dc.SetColor(color.Black)
		// Разбиваем длинные заголовки на несколько строк
		if len(col.Title) > 10 && strings.Contains(col.Title, " ") {
			parts := strings.Split(col.Title, " ")
			if len(parts) == 2 {
				dc.DrawStringAnchored(parts[0], currentX+col.Width/2, y+rowHeight/2-7, 0.5, 0.5)
				dc.DrawStringAnchored(parts[1], currentX+col.Width/2, y+rowHeight/2+7, 0.5, 0.5)
				currentX += col.Width
				continue
			}
		}
		dc.DrawStringAnchored(col.Title, currentX+col.Width/2, y+rowHeight/2, 0.5, 0.5)
		currentX += col.Width
	}

	// Данные
	for _, report := range reports {
		y += rowHeight
		currentX = margin
		dc.SetColor(color.Black)

		// Форматируем данные
		values := []string{
			report.Ticker,
			report.Name,
			report.Currencies,
			formatInt(report.Quantity),
			formatPercent(report.PercentOfPortfolio),
			formatTime(report.MaturityDate),
			formatInt(report.Duration),
			formatTime(report.BuyDate),
			formatFloat(report.PositionPrice),
			formatPercent(report.YieldToMaturityOnPurchase),
			formatPercent(report.YieldToMaturity),
			formatFloat(report.CurrentPrice),
			formatFloat(report.Nominal),
			formatFloat(report.Profit),
			formatPercent(report.ProfitInPercentage),
		}

		for i, col := range columns {
			dc.DrawStringAnchored(
				values[i],
				currentX+col.Width/2,
				y+rowHeight/2,
				0.5, 0.5,
			)
			currentX += col.Width
		}
	}

	return dc.SavePNG(filename)

}

func GenerateTablePNG(reports []service_models.GeneralBondReportPosition, typeOfBonds string) (_ []byte, err error) {
	defer func() { err = e.WrapIfErr("can't generate table in png", err) }()
	const (
		width        = 1200 // Увеличим ширину для лучшего отображения
		heightPerRow = 40
		margin       = 20.0
		headerHeight = 60
		rowHeight    = 30
		colCount     = 15 // Количество отображаемых колонок
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
		return nil, err
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(face)

	// Заливаем фон
	dc.SetColor(color.White)
	dc.Clear()

	// Рисуем заголовок
	dc.SetColor(color.Black)
	var name string
	switch typeOfBonds {
	case service_models.ReplacedBonds:
		name = "Отчет по замещающим облигациям"
	case service_models.RubBonds:
		name = "Отчет по рублевым облигациям"
	case service_models.EuroBonds:
		name = "Отчет по валютным облигациям"
	}
	dc.DrawStringAnchored(name, width/2, margin, 0.5, 0.5)
	dc.SetFontFace(face)

	// Определяем колонки для отображения с более подходящими ширинами
	columns := []struct {
		Title string
		Width float64
		Flex  float64
	}{
		{"Тикер", 120, 0},              // 1
		{"Название", 180, 2},           // 2
		{"Валюта", 60, 0},              // 3
		{"Кол-во", 70, 0},              // 4
		{"% от портфеля", 90, 0},       // 5
		{"Дата погашения", 110, 0},     // 6
		{"Дюрация", 80, 0},             // 7
		{"Дата покупки", 90, 0},        // 8
		{"Ср.цена", 80, 0},             // 9
		{"Дох-ть при покупке", 120, 0}, // 10
		{"Дох-ть тек.", 90, 0},         // 11 (сокращенный заголовок)
		{"Тек.цена", 80, 0},            // 12 (сокращенный заголовок)
		{"Номинал", 80, 0},             // 13
		{"Доход", 80, 0},               // 14
		{"Дох-ть в %", 80, 0},          // 15
	}

	// Проверяем, что сумма минимальных ширин не превышает доступную ширину
	totalMinWidth := 0.0
	for _, col := range columns {
		totalMinWidth += col.Width
	}

	availableWidth := float64(width) - 2*margin
	if totalMinWidth > availableWidth {
		// Если минимальные ширины не помещаются, масштабируем их
		scaleFactor := availableWidth / totalMinWidth
		for i := range columns {
			columns[i].Width *= scaleFactor
		}
	} else {
		// Распределяем оставшееся пространство пропорционально Flex
		remainingWidth := availableWidth - totalMinWidth
		totalFlex := 0.0
		for _, col := range columns {
			totalFlex += col.Flex
		}

		if totalFlex > 0 {
			flexUnit := remainingWidth / totalFlex
			for i := range columns {
				if columns[i].Flex > 0 {
					columns[i].Width += flexUnit * columns[i].Flex
				}
			}
		}
	}

	// Рисование таблицы
	y := float64(headerHeight)
	currentX := margin

	// Заголовки
	for _, col := range columns {
		dc.SetColor(color.RGBA{200, 200, 200, 255})
		dc.DrawRectangle(currentX, y, col.Width, rowHeight)
		dc.Fill()

		dc.SetColor(color.Black)
		// Разбиваем длинные заголовки на несколько строк
		if len(col.Title) > 10 && strings.Contains(col.Title, " ") {
			parts := strings.Split(col.Title, " ")
			if len(parts) == 2 {
				dc.DrawStringAnchored(parts[0], currentX+col.Width/2, y+rowHeight/2-7, 0.5, 0.5)
				dc.DrawStringAnchored(parts[1], currentX+col.Width/2, y+rowHeight/2+7, 0.5, 0.5)
				currentX += col.Width
				continue
			}
		}
		dc.DrawStringAnchored(col.Title, currentX+col.Width/2, y+rowHeight/2, 0.5, 0.5)
		currentX += col.Width
	}

	// Данные
	for _, report := range reports {
		y += rowHeight
		currentX = margin
		dc.SetColor(color.Black)

		// Форматируем данные
		values := []string{
			report.Ticker,
			report.Name,
			report.Currencies,
			formatInt(report.Quantity),
			formatPercent(report.PercentOfPortfolio),
			formatTime(report.MaturityDate),
			formatInt(report.Duration),
			formatTime(report.BuyDate),
			formatFloat(report.PositionPrice),
			formatPercent(report.YieldToMaturityOnPurchase),
			formatPercent(report.YieldToMaturity),
			formatFloat(report.CurrentPrice),
			formatFloat(report.Nominal),
			formatFloat(report.Profit),
			formatPercent(report.ProfitInPercentage),
		}

		for i, col := range columns {
			dc.DrawStringAnchored(
				values[i],
				currentX+col.Width/2,
				y+rowHeight/2,
				0.5, 0.5,
			)
			currentX += col.Width
		}
	}
	pngData, err := EncodePNGToBuffer(dc)
	if err != nil {
		return nil, err
	}

	return pngData, nil

}

func EncodePNGToBuffer(dc *gg.Context) ([]byte, error) {
	img := dc.Image()
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func formatInt(value int64) string {
	return fmt.Sprintf("%v", value)
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}
func formatTime(value time.Time) string {
	return value.Format(layout)
}
