package endpoints

import (
	"microblog-app/internal/session"
	"net/http"
)

func Logout() Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		if _, err := session.Get(r, w); err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}
		if err := session.End(r, w); err != nil {
			w.WriteHeader(500)
			return err
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}
}
