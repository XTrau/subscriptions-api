package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (srw *responseWriter) WriteHeader(statusCode int) {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{w, 200}

			next.ServeHTTP(rw, r)

			elapsed := time.Since(start)

			logger.Info(
				"request",
				slog.String("Method", r.Method),
				slog.String("Path", r.URL.Path),
				slog.Int("Status", rw.statusCode),
				slog.Duration("Duration", elapsed),
			)
		})
	}

}
