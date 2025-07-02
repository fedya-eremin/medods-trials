package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"example.com/m/internal/contextkeys"
)

type responseWriterWithLog struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWithLog) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func LoggerMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLogger := logger.With(
			"method", r.Method,
			"path", r.URL.Path,
			"remoteAddr", r.RemoteAddr,
		)
		ctx := contextkeys.WithContextValue(r.Context(), contextkeys.LoggerKey, logger)
		r = r.WithContext(ctx)

		ww := &responseWriterWithLog{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		requestLogger = requestLogger.With(
			"status", ww.statusCode,
			"duration", duration,
		)
		level := slog.LevelInfo
		if ww.statusCode >= 500 {
			level = slog.LevelError
		} else if ww.statusCode >= 400 {
			level = slog.LevelWarn
		}
		requestLogger.LogAttrs(ctx, level, "request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.Int("status", ww.statusCode),
			slog.String("duration", duration.String()),
			slog.Float64("duration_ms", float64(duration.Microseconds())/1000),
		)
	})
}
