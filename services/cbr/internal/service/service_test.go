//go:build unit

package service

import (
	"cbr/internal/clients/cbr/mocks"
	"cbr/internal/models"
	"cbr/internal/utils"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	layout  = "02/01/2006"
	cbrHost = "www.cbr.ru"
)

func TestGetAllCurrencies(t *testing.T) {
	ctx := context.Background()
	logg := slog.New(slog.NewTextHandler(io.Discard, nil))
	now := time.Now()
	location, err := utils.GetMoscowLocation()
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		wantMock := models.CurrenciesResponce{Date: "date"}
		cbrClientMock := mocks.NewCbrClient(t)
		cbrClientMock.On("GetAllCurrencies", ctx, mock.AnythingOfType("string")).
			Return(wantMock, nil).Once()
		srv := NewService(logg, cbrClientMock, location)
		currResp, err := srv.GetAllCurrencies(ctx, now)
		require.NoError(t, err)
		require.Equal(t, wantMock, currResp)
	})
	t.Run("GetAllCurrencies error", func(t *testing.T) {
		errContains := "failed to get all currencies from client"
		cbrClientMock := mocks.NewCbrClient(t)
		cbrClientMock.On("GetAllCurrencies", ctx, mock.AnythingOfType("string")).
			Return(models.CurrenciesResponce{}, errors.New("could not do request")).Once()
		srv := NewService(logg, cbrClientMock, location)
		currResp, err := srv.GetAllCurrencies(ctx, now)
		require.Error(t, err)
		require.ErrorContains(t, err, errContains)
		require.Empty(t, currResp)
	})
}
