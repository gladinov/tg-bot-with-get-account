//go:build unit

package handlers

import (
	"cbr/internal/service/mocks"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	httpheaders "github.com/gladinov/contracts/http"
	"github.com/gladinov/valuefromcontext"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestContextHeaderTraceIdMiddleware_FromHeader(t *testing.T) {
	e := echo.New()

	srvc := mocks.NewCurrencyService(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewHandlers(logger, srvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(httpheaders.HeaderTraceID, "trace-123")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	next := func(c echo.Context) error {
		called = true

		traceID, err := valuefromcontext.GetTraceID(c.Request().Context())
		require.NoError(t, err)
		require.Equal(t, "trace-123", traceID)

		return nil
	}

	mw := h.ContextHeaderTraceIdMiddleWare(next)
	err := mw(c)

	require.NoError(t, err)
	require.True(t, called)
}

func TestContextHeaderTraceIdMiddleware_EmptyTraceID(t *testing.T) {
	e := echo.New()

	srvc := mocks.NewCurrencyService(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewHandlers(logger, srvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(httpheaders.HeaderTraceID, "")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	next := func(c echo.Context) error {
		called = true

		traceID, err := valuefromcontext.GetTraceID(c.Request().Context())
		require.NoError(t, err)
		require.NotEmpty(t, traceID)

		return nil
	}

	mw := h.ContextHeaderTraceIdMiddleWare(next)
	err := mw(c)

	require.NoError(t, err)
	require.True(t, called)
}

func TestLoggerMiddleware_PassesThrough(t *testing.T) {
	e := echo.New()

	srvc := mocks.NewCurrencyService(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewHandlers(logger, srvc)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	nextCalled := false
	next := func(c echo.Context) error {
		nextCalled = true
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}

	mw := h.LoggerMiddleWare(next)
	err := mw(c)

	require.NoError(t, err)
	require.True(t, nextCalled)
}

func TestLoggerMiddleware_Error(t *testing.T) {
	e := echo.New()

	srvc := mocks.NewCurrencyService(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewHandlers(logger, srvc)

	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	next := func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	mw := h.LoggerMiddleWare(next)
	err := mw(c)

	require.Error(t, err)
}
