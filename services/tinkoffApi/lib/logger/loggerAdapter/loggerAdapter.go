package loggeradapter

import (
	"fmt"
	"log/slog"
)

type LoggerAdapter struct {
	slogLogger *slog.Logger
}

func NewLoggerAdapter(slogLogger *slog.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		slogLogger: slogLogger,
	}
}

func (l *LoggerAdapter) Infof(template string, args ...any) {
	slogAttr := make([]any, 0, len(args))
	for i, arg := range args {
		count := i + 1
		key := fmt.Sprintf("attr_%v", count)
		slogAttr = append(slogAttr, slog.Any(
			key,
			arg))
	}
	l.slogLogger.Info(template, slogAttr...)
}

func (l *LoggerAdapter) Errorf(template string, args ...any) {
	slogAttr := make([]any, 0, len(args))
	for i, arg := range args {
		count := i + 1
		key := fmt.Sprintf("attr_%v", count)
		slogAttr = append(slogAttr, slog.Any(
			key,
			arg))
	}
	l.slogLogger.Error(template, slogAttr...)
}
func (l *LoggerAdapter) Fatalf(template string, args ...any) {
	slogAttr := make([]any, 0, len(args))
	for i, arg := range args {
		count := i + 1
		key := fmt.Sprintf("attr_%v", count)
		slogAttr = append(slogAttr, slog.Any(
			key,
			arg))
	}
	l.slogLogger.Error(template, slogAttr...)
}
