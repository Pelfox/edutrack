package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	loggerContextKey = "logger"
	requestIDHeader  = "X-Request-ID"
)

// RequestLoggerMiddleware добавляет к запросу идентификатор и пишет краткий итог обработки.
func RequestLoggerMiddleware(baseLogger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx.Header(requestIDHeader, requestID)
		requestLogger := baseLogger.With().
			Str("request_id", requestID).
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("client_ip", ctx.ClientIP()).
			Logger()

		ctx.Set(loggerContextKey, requestLogger)

		startedAt := time.Now()
		ctx.Next()

		status := ctx.Writer.Status()
		event := requestLogger.Info()
		if status >= http.StatusInternalServerError {
			event = requestLogger.Error()
		} else if status >= http.StatusBadRequest {
			event = requestLogger.Warn()
		}

		event.
			Int("status", status).
			Dur("latency", time.Since(startedAt)).
			Int("response_size", ctx.Writer.Size()).
			Str("route", requestRoute(ctx)).
			Msg("request completed")
	}
}

// RecoveryMiddleware перехватывает panic и пишет его через общий логгер приложения.
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, recovered any) {
		logger := loggerFromContext(ctx)
		logger.
			Error().
			Interface("panic", recovered).
			Str("route", requestRoute(ctx)).
			Msg("panic recovered")

		ctx.AbortWithStatus(http.StatusInternalServerError)
	})
}

func loggerFromContext(ctx *gin.Context) zerolog.Logger {
	value, ok := ctx.Get(loggerContextKey)
	if !ok {
		return zerolog.Nop()
	}

	logger, ok := value.(zerolog.Logger)
	if !ok {
		return zerolog.Nop()
	}

	return logger
}

func requestRoute(ctx *gin.Context) string {
	route := ctx.FullPath()
	if route == "" {
		return ctx.Request.URL.Path
	}

	return route
}
