package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/mesamhaider/edge-device-sample/internal/logging"
)

func RequestContextMiddleware(baseLogger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = logging.RequestIDFromContext(r.Context())
			}
			if requestID == "" {
				requestID = logging.GenerateRequestID()
			}

			ctx := logging.ContextWithRequestID(r.Context(), requestID)
			requestLogger := baseLogger.With(
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
			ctx = logging.ContextWithLogger(ctx, requestLogger)

			w.Header().Set("X-Request-ID", requestID)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			requestLogger.Info("request started")

			next.ServeHTTP(ww, r.WithContext(ctx))

			duration := time.Since(start)
			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			requestLogger.Info("request completed",
				zap.Int("status", status),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", duration),
			)
		})
	}
}
