//go:build unit

package service

import (
	"bonds-report-service/internal/application/mocks"
	factories "bonds-report-service/internal/application/testing"
	"bonds-report-service/internal/domain"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getServiceForTestMoex(logger *slog.Logger, externalApis *ExternalApis) *Service {
	return NewService(logger, nil, externalApis, nil, nil, nil, nil, nil)
}

func TestGetSpecificationsFromMoex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	fixedTime := time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		ticker := "test_ticker"
		want := factories.NewValuesMoex()
		mockMoex := mocks.NewMoexClient(t)
		externalApis := NewExternalApis(mockMoex, nil, nil)
		srvs := getServiceForTestMoex(logger, externalApis)
		mockMoex.On("GetSpecifications", ctx, ticker, fixedTime).
			Return(want, nil)

		got, err := srvs.GetSpecificationsFromMoex(ctx, ticker, fixedTime)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("Err: Empty Ticker", func(t *testing.T) {
		ticker := ""

		mockMoex := mocks.NewMoexClient(t)
		externalApis := NewExternalApis(mockMoex, nil, nil)
		srvs := getServiceForTestMoex(logger, externalApis)

		_, err := srvs.GetSpecificationsFromMoex(ctx, ticker, fixedTime)
		assert.Error(t, err)
	})
	t.Run("Err: GetSpecifications", func(t *testing.T) {
		ticker := "test_ticker"
		wanterr := "failed get spec"
		mockMoex := mocks.NewMoexClient(t)
		externalApis := NewExternalApis(mockMoex, nil, nil)
		srvs := getServiceForTestMoex(logger, externalApis)
		mockMoex.On("GetSpecifications", ctx, ticker, fixedTime).
			Return(domain.ValuesMoex{}, errors.New("failed get spec"))

		_, err := srvs.GetSpecificationsFromMoex(ctx, ticker, fixedTime)
		assert.Error(t, err)
		assert.ErrorContains(t, err, wanterr)
	})
}
