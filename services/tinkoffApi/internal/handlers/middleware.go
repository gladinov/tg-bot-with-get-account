package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"tinkoffApi/lib/cryptoToken"
	"tinkoffApi/lib/valuefromcontext"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func (h *Handlers) AuthCheckTokenMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		const op = "handlers.AuthCheckTokenMiddleWare"
		logg := h.logger.With(slog.String("op", op))
		logg.Debug("start", slog.String("path", c.Path()))

		defer func() {
			if err != nil {
				logg.Error("auth middleware failed", slog.Any("error", err))
			} else {
				logg.Info("success finished")
			}
		}()

		if c.Path() == "/tinkoff/checktoken" {
			return next(c)
		}

		chatID := c.Request().Header.Get(HeaderChatID)
		if chatID == "" {
			err = errHeaderRequired
			return echo.NewHTTPError(http.StatusUnauthorized, errHeaderRequired)
		}
		ctx := c.Request().Context()
		tokenInBase64, err := h.redis.Get(ctx, chatID).Result()
		switch err {
		case nil:
		case redis.Nil:
			return echo.NewHTTPError(http.StatusUnauthorized, errNoTokenInRedis)
		default:
			return echo.NewHTTPError(http.StatusServiceUnavailable, errRedisDoNotAnswer)
		}

		encryptedToken, err := cryptoToken.GetEncryptedTokenFromBase64(tokenInBase64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errHeaderRequired)
		}
		token, err := cryptoToken.DecryptToken(&encryptedToken, h.tokenCrypter.Key)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errInvalidAuthFormat)
		}
		ctx = context.WithValue(c.Request().Context(), valuefromcontext.EncryptedTokenKey, token)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (h *Handlers) AuthCheckTokenInHeadersMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		const op = "handlers.AuthCheckTokenInHeadersMiddleWare"
		logg := h.logger.With(slog.String("op", op))
		logg.Debug("start")
		defer func() {
			if err != nil {
				logg.Warn("auth failed",
					slog.Any("error", err),
					slog.String("path", c.Path()),
					slog.String("method", c.Request().Method))
			} else {
				logg.Debug("auth success",
					slog.String("path", c.Path()),
					slog.String("method", c.Request().Method))
			}
		}()

		tokenInBase64 := c.Request().Header.Get(HeaderEncryptedToken)
		if tokenInBase64 == "" {
			err = errHeaderRequired
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		encryptedToken, err := cryptoToken.GetEncryptedTokenFromBase64(tokenInBase64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errHeaderRequired)
		}
		token, err := cryptoToken.DecryptToken(&encryptedToken, h.tokenCrypter.Key)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errInvalidAuthFormat)
		}
		ctx := context.WithValue(c.Request().Context(), valuefromcontext.EncryptedTokenKey, token)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

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
