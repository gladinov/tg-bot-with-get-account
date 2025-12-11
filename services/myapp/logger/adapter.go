package loggAdapter

import (
	"fmt"
	"log/slog"
	"os"
)

type Adapter struct {
	slog *slog.Logger
}

func New(logger *slog.Logger) *Adapter {
	return &Adapter{slog: logger}
}

func (s *Adapter) Infof(template string, args ...any) {
	s.slog.Info(fmt.Sprintf(template, args...))
}

func (s *Adapter) Errorf(template string, args ...any) {
	s.slog.Error(fmt.Sprintf(template, args...))
}

func (s *Adapter) Printf(template string, args ...any) {
	s.slog.Info(fmt.Sprintf(template, args...))
}

func (s *Adapter) Fatalf(template string, args ...any) {
	s.slog.Error(fmt.Sprintf(template, args...))
	os.Exit(1)
}
