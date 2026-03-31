//go:build unit

package handlers

import (
	"errors"
	"io"
	"log/slog"
	"moex/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHTTPErrorHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := HTTPErrorHandler(logger)

	t.Run("Internal Server Error logs", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := errors.New("something went wrong")
		handler(err, c)

		require.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp models.ErrorResponse
		require.NoError(t, c.JSON(rec.Code, resp))
	})

	t.Run("HTTPError < 500 does not log", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		he := echo.NewHTTPError(http.StatusBadRequest, "bad request")
		handler(he, c)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var resp models.ErrorResponse
		require.NoError(t, c.JSON(rec.Code, resp))
	})

	t.Run("HTTPError >= 500 logs", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		he := echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		handler(he, c)

		require.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp models.ErrorResponse
		require.NoError(t, c.JSON(rec.Code, resp))
	})
}
