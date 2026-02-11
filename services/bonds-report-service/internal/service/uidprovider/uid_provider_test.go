package uidprovider

import (
	"bonds-report-service/internal/models/domain"
	repoMock "bonds-report-service/internal/repository/mocks"
	serviceMock "bonds-report-service/internal/service/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUidProvider_GetUid(t *testing.T) {
	ctx := context.Background()
	const instrumentUid = "instr1"
	const expectedUid = "asset123"

	t.Run("success from storage", func(t *testing.T) {
		storageMock := repoMock.NewStorage(t)
		analyticsMock := serviceMock.NewTinkoffAnalyticsClient(t)

		// Устанавливаем время фиктивное
		now := time.Now()
		provider := NewUidProvider(storageMock, analyticsMock)
		provider.now = func() time.Time { return now }

		// storage возвращает дату обновления недавно, UID есть
		storageMock.On("IsUpdatedUids", ctx).Return(now, nil)
		storageMock.On("GetUid", ctx, instrumentUid).Return(expectedUid, nil)

		uid, err := provider.GetUid(ctx, instrumentUid)
		assert.NoError(t, err)
		assert.Equal(t, expectedUid, uid)

		storageMock.AssertExpectations(t)
	})

	t.Run("TTL expired triggers UpdateAndGetUid", func(t *testing.T) {
		storageMock := repoMock.NewStorage(t)
		analyticsMock := serviceMock.NewTinkoffAnalyticsClient(t)

		past := time.Now().Add(-2 * HoursToUpdate)
		provider := NewUidProvider(storageMock, analyticsMock)
		provider.now = func() time.Time { return time.Now() }

		storageMock.On("IsUpdatedUids", ctx).Return(past, nil)
		analyticsMock.On("GetAllAssetUids", ctx).Return(map[string]string{
			instrumentUid: expectedUid,
		}, nil)
		storageMock.On("SaveUids", ctx, mock.Anything).Return(nil)

		uid, err := provider.GetUid(ctx, instrumentUid)
		assert.NoError(t, err)
		assert.Equal(t, expectedUid, uid)

		storageMock.AssertExpectations(t)
		analyticsMock.AssertExpectations(t)
	})

	t.Run("storage returns ErrEmptyUids triggers UpdateAndGetUid", func(t *testing.T) {
		storageMock := repoMock.NewStorage(t)
		analyticsMock := serviceMock.NewTinkoffAnalyticsClient(t)

		provider := NewUidProvider(storageMock, analyticsMock)

		storageMock.On("IsUpdatedUids", ctx).Return(time.Time{}, domain.ErrEmptyUids)
		analyticsMock.On("GetAllAssetUids", ctx).Return(map[string]string{
			instrumentUid: expectedUid,
		}, nil)
		storageMock.On("SaveUids", ctx, mock.Anything).Return(nil)

		uid, err := provider.GetUid(ctx, instrumentUid)
		assert.NoError(t, err)
		assert.Equal(t, expectedUid, uid)

		storageMock.AssertExpectations(t)
		analyticsMock.AssertExpectations(t)
	})

	t.Run("UpdateAndGetUid returns ErrEmptyUidAfterUpdate", func(t *testing.T) {
		storageMock := repoMock.NewStorage(t)
		analyticsMock := serviceMock.NewTinkoffAnalyticsClient(t)

		provider := NewUidProvider(storageMock, analyticsMock)
		storageMock.On("IsUpdatedUids", ctx).Return(time.Time{}, domain.ErrEmptyUids)
		analyticsMock.On("GetAllAssetUids", ctx).Return(map[string]string{}, nil)

		uid, err := provider.GetUid(ctx, instrumentUid)
		assert.ErrorIs(t, err, domain.ErrEmptyUidAfterUpdate)
		assert.Equal(t, "", uid)
	})
	t.Run("IsUpdatedUids returns unexpected error", func(t *testing.T) {
		storageMock := repoMock.NewStorage(t)
		analyticsMock := serviceMock.NewTinkoffAnalyticsClient(t)

		provider := NewUidProvider(storageMock, analyticsMock)

		storageMock.On("IsUpdatedUids", ctx).Return(time.Time{}, errors.New("db failure"))

		uid, err := provider.GetUid(ctx, instrumentUid)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "check updated uids")
		assert.Empty(t, uid)

		storageMock.AssertExpectations(t)
	})

	t.Run("UID missing in storage after recent update", func(t *testing.T) {
		storageMock := repoMock.NewStorage(t)
		analyticsMock := serviceMock.NewTinkoffAnalyticsClient(t)

		now := time.Now()
		provider := NewUidProvider(storageMock, analyticsMock)
		provider.now = func() time.Time { return now }

		// Storage говорит, что обновление было недавно
		storageMock.On("IsUpdatedUids", ctx).Return(now, nil)
		// Но UID отсутствует
		storageMock.On("GetUid", ctx, instrumentUid).Return("", domain.ErrEmptyUids)

		uid, err := provider.GetUid(ctx, instrumentUid)
		assert.ErrorIs(t, err, domain.ErrEmptyUidAfterUpdate)
		assert.Empty(t, uid)

		storageMock.AssertExpectations(t)
	})
}

func TestUidProvider_UpdateAndGetUid(t *testing.T) {
	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		mockAnalytics := serviceMock.NewTinkoffAnalyticsClient(t)
		mockStorage := repoMock.NewStorage(t)

		provider := NewUidProvider(mockStorage, mockAnalytics)
		instrUid := "instr123"
		expectedUid := "asset456"

		mockAnalytics.On("GetAllAssetUids", ctx).Return(map[string]string{
			instrUid: expectedUid,
		}, nil)

		mockStorage.On("SaveUids", ctx, mock.Anything).Return(nil)

		uid, err := provider.UpdateAndGetUid(ctx, instrUid)
		assert.NoError(t, err)
		assert.Equal(t, expectedUid, uid)

		mockAnalytics.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("analytics error", func(t *testing.T) {
		mockAnalytics := serviceMock.NewTinkoffAnalyticsClient(t)
		mockStorage := repoMock.NewStorage(t)

		provider := NewUidProvider(mockStorage, mockAnalytics)
		instrUid := "instr123"
		mockAnalytics.On("GetAllAssetUids", ctx).Return(nil, errors.New("network error"))

		uid, err := provider.UpdateAndGetUid(ctx, instrUid)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get all asset uids")
		assert.Empty(t, uid)

		mockAnalytics.AssertExpectations(t)
	})

	t.Run("uid not found", func(t *testing.T) {
		mockAnalytics := serviceMock.NewTinkoffAnalyticsClient(t)
		mockStorage := repoMock.NewStorage(t)

		provider := NewUidProvider(mockStorage, mockAnalytics)
		instrUid := "instr123"
		mockAnalytics.On("GetAllAssetUids", ctx).Return(map[string]string{}, nil)

		uid, err := provider.UpdateAndGetUid(ctx, instrUid)
		assert.ErrorIs(t, err, domain.ErrEmptyUidAfterUpdate)
		assert.Empty(t, uid)

		mockAnalytics.AssertExpectations(t)
	})

	t.Run("storage save error", func(t *testing.T) {
		mockAnalytics := serviceMock.NewTinkoffAnalyticsClient(t)
		mockStorage := repoMock.NewStorage(t)

		provider := NewUidProvider(mockStorage, mockAnalytics)
		instrUid := "instr123"
		expectedUid := "asset456"

		mockAnalytics.On("GetAllAssetUids", ctx).Return(map[string]string{
			instrUid: expectedUid,
		}, nil)
		mockStorage.On("SaveUids", ctx, mock.Anything).Return(errors.New("db write failed"))

		uid, err := provider.UpdateAndGetUid(ctx, instrUid)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "save uids")
		assert.Empty(t, uid)

		mockAnalytics.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}
