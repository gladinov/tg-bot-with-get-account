package handlers

import (
	"bonds-report-service/lib/valuefromcontext"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.AuthMiddleware"
		chatID := c.GetHeader(valuefromcontext.HeaderChatID)
		if chatID == "" {
			logger.Warn("missing chat id header",
				slog.String("op", op),
				slog.String("header", valuefromcontext.HeaderChatID),
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing X-Chat-ID header"})
			return
		}
		ctx := context.WithValue(c.Request.Context(), valuefromcontext.ChatIdKey, chatID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(
			slog.String("component", "middleware/logger"),
		)

		req := c.Request
		resp := c.Writer
		entry := logg.With(
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.String("remote_addr", req.RemoteAddr),
			slog.String("user_agent", req.UserAgent()),
			// slog.String("request_id", req.Header.Get(models.RequestIDHeader)),
		)
		t1 := time.Now()
		c.Next()
		duration := time.Since(t1).Milliseconds()
		entry.Info("request completed",
			slog.Int("status", resp.Status()),
			slog.Int64("bytes", int64(resp.Size())),
			slog.Int64("duration", duration),
		)

	}
}
