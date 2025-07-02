package contextkeys

import (
	"context"
	"log/slog"
)

type ContextKey string

const (
	LoggerKey    ContextKey = "masterLogger"
	JWTClaimsKey ContextKey = "jwtClaims"
)

func WithContextValue(ctx context.Context, key ContextKey, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetContextValue[V any](ctx context.Context, key ContextKey) (V, bool) {
	val, ok := ctx.Value(key).(V)
	return val, ok
}

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := GetContextValue[*slog.Logger](ctx, LoggerKey); ok {
		return logger
	}
	return slog.Default()
}
