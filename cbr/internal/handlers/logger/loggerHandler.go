package loggerhandler

import (
	"cbr/lib/models"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			log := log.With(
				slog.String("component", "middleware/logger"),
			)

			//log.Info("logger middleware enabled")

			req := c.Request()
			resp := c.Response()

			entry := log.With(
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("remote_addr", req.RemoteAddr),
				slog.String("user_agent", req.UserAgent()),
				slog.String("request_id", req.Header.Get(models.RequestIDHeader)),
			)
			t1 := time.Now()
			err := next(c)
			duration := time.Since(t1).Milliseconds()
			entry.Info("request completed",
				slog.Int("status", resp.Status),
				slog.Int64("bytes", resp.Size),
				slog.Int64("duration", duration),
			)

			return err
		}
	}
}
