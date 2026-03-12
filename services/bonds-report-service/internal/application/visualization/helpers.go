package visualization

import (
	"fmt"
	"time"
)

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
