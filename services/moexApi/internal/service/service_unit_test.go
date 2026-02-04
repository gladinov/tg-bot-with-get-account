//go:build unit

package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"moex/internal/clients/moex/mocks"
	"moex/internal/models"
	"moex/internal/testdata/factories"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetSpecifications(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))
	date := time.Date(2025, 11, 16, 0, 0, 0, 0, time.UTC)
	ticker := "OFZ26238"

	t.Run("Success", func(t *testing.T) {
		mockMoexClient := mocks.NewMoexClient(t)
		req := models.SpecificationsRequest{
			Ticker: ticker,
			Date:   date,
		}
		wantMock := factories.NewSpecificationsResponse()
		want := factories.NewValues()
		mockMoexClient.On("GetSpecifications",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("time.Time")).
			Return(wantMock, nil).Once()
		serviceClient := NewServiceClient(logg, mockMoexClient)
		got, err := serviceClient.GetSpecifications(ctx, req)
		require.NoError(t, err)
		require.Equal(t, want, got)

		mockMoexClient.AssertExpectations(t)
	})
	t.Run("Err:Get specification err", func(t *testing.T) {
		mockMoexClient := mocks.NewMoexClient(t)
		req := models.SpecificationsRequest{
			Ticker: ticker,
			Date:   date,
		}
		wantMock := models.SpecificationsResponce{}
		want := models.Values{}
		errContains := "could not get specification from moexClient"
		mockMoexClient.On("GetSpecifications",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("time.Time")).
			Return(wantMock, errors.New("GetSpecification failed")).Once()
		serviceClient := NewServiceClient(logg, mockMoexClient)
		got, err := serviceClient.GetSpecifications(ctx, req)
		require.Error(t, err)
		require.ErrorContains(t, err, errContains)
		require.Equal(t, want, got)

		mockMoexClient.AssertExpectations(t)
	})
	t.Run("Err:No data in MOEX", func(t *testing.T) {
		mockMoexClient := mocks.NewMoexClient(t)
		timesCallMock := 14
		req := models.SpecificationsRequest{
			Ticker: ticker,
			Date:   date,
		}
		wantMock := models.SpecificationsResponce{}
		want := models.Values{}
		errContains := "could not find data in MOEX"
		mockMoexClient.On("GetSpecifications",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("time.Time")).
			Return(wantMock, nil).Times(timesCallMock)
		serviceClient := NewServiceClient(logg, mockMoexClient)
		got, err := serviceClient.GetSpecifications(ctx, req)
		require.Error(t, err)
		require.ErrorContains(t, err, errContains)
		require.Equal(t, want, got)

		mockMoexClient.AssertExpectations(t)
	})
}

func TestClampDate_FutureDate(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	future := now.Add(24 * time.Hour)

	got := clampDate(future, now)

	require.Equal(t, now, got)
}
