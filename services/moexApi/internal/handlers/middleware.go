package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) LoggerMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		logg := h.logger.With(
			slog.String("component", "middleware/logger"),
		)

		req := c.Request()
		resp := c.Response()
		entry := logg.With(
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.String("remote_addr", req.RemoteAddr),
			slog.String("user_agent", req.UserAgent()),
			// slog.String("request_id", req.Header.Get(models.RequestIDHeader)),
		)
		start := time.Now()

		defer func() {
			status := resp.Status
			if status == 0 {
				status = http.StatusOK
			}

			attrs := []any{
				slog.Int("status", status),
				slog.Int64("bytes", resp.Size),
				slog.Duration("duration", time.Since(start)),
			}

			if err != nil {
				entry.Error("request failed",
					append(attrs, slog.Any("error", err))...,
				)
			} else {
				entry.Info("request completed", attrs...)
			}
		}()
		err = next(c)

		return err
	}
}
