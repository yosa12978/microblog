package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Recovery(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(
						"panic recovery",
						"method", r.Method,
						"path", r.URL.Path,
						"error", fmt.Sprintf("%+v", err),
					)
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
