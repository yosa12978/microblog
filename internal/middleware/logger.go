package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
			start := time.Now().UTC()
			next.ServeHTTP(writer, r)
			latency := time.Since(start).Microseconds()
			logger.Info(
				"incoming request",
				"latency_us", latency,
				"status_code", writer.statusCode,
				"method", r.Method,
				//"user_agent", r.UserAgent(),
				"path", r.URL.Path,
			)
		})
	}
}
