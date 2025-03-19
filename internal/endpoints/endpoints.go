package endpoints

import "net/http"

type Endpoint func(w http.ResponseWriter, r *http.Request) error

func (f Endpoint) Unwrap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// render error template
			w.Write([]byte(err.Error()))
			return
		}
	}
}
