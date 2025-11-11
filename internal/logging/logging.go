package logging

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ctxKey string

const (
	loggerKey    ctxKey = "logger"
	requestIDKey ctxKey = "request_id"
)

var noopLogger = zap.NewNop()

// ContextWithLogger stores the provided logger in the context.
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext retrieves a logger from context, falling back to the provided default.
func LoggerFromContext(ctx context.Context, fallback *zap.Logger) *zap.Logger {
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok && ctxLogger != nil {
		return ctxLogger
	}
	if fallback != nil {
		return fallback
	}
	return noopLogger
}

// ContextWithRequestID stores the request ID in the context.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext retrieves the request ID from the context.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// GenerateRequestID returns a random hex-encoded 16-byte identifier.
func GenerateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
