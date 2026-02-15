//go:build unit

package bondreport

import (
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/service/bondReport/mocks"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	factories "bonds-report-service/internal/service/testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSpecificationsFromMoex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	fixedTime := time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		ticker := "test_ticker"
		want := factories.NewValuesMoex()
		mockMoex := mocks.NewMoexClient(t)
		bondReporter := NewBondReporter(logger, mockMoex)
		mockMoex.On("GetSpecifications", ctx, ticker, fixedTime).
			Return(want, nil)

		got, err := bondReporter.GetSpecificationsFromMoex(ctx, ticker, fixedTime)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("Err: Empty Ticker", func(t *testing.T) {
		ticker := ""

		mockMoex := mocks.NewMoexClient(t)
		bondReporter := NewBondReporter(logger, mockMoex)

		_, err := bondReporter.GetSpecificationsFromMoex(ctx, ticker, fixedTime)
		assert.Error(t, err)
	})
	t.Run("Err: GetSpecifications", func(t *testing.T) {
		ticker := "test_ticker"
		wanterr := "failed get spec"
		mockMoex := mocks.NewMoexClient(t)
		bondReporter := NewBondReporter(logger, mockMoex)
		mockMoex.On("GetSpecifications", ctx, ticker, fixedTime).
			Return(domain.ValuesMoex{}, errors.New("failed get spec"))

		_, err := bondReporter.GetSpecificationsFromMoex(ctx, ticker, fixedTime)
		assert.Error(t, err)
		assert.ErrorContains(t, err, wanterr)
	})
}
