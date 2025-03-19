package middleware

import (
	"microblog-app/internal/session"
	"net/http"
)

func Anon() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := session.Get(r, w); err == nil {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Forbidden", 403)
		})
	}
}

func Admin() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := session.Get(r, w); err == nil {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Forbidden", 403)
		})
	}
}
