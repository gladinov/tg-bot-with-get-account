package loggerdicard

type LoggerDiscard struct{}

func NewLoggerDiscard() *LoggerDiscard {
	return &LoggerDiscard{}
}

func (l *LoggerDiscard) Infof(template string, args ...any)  {}
func (l *LoggerDiscard) Errorf(template string, args ...any) {}
func (l *LoggerDiscard) Fatalf(template string, args ...any) {}
