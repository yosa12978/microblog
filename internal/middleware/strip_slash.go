package middleware

import "net/http"

func StripSlash() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == "" {
				path = "/"
			}
			if len(path) > 1 && path[len(path)-1] == '/' {
				path = path[:len(path)-1]
			}
			r.URL.Path = path
			next.ServeHTTP(w, r)
		})
	}
}
