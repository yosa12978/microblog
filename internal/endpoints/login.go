package endpoints

import (
	"microblog-app/internal/config"
	"microblog-app/internal/dto"
	"microblog-app/internal/session"
	"microblog-app/internal/templates"
	"microblog-app/pkg"
	"net/http"
)

func LoginGet() Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		if _, err := session.Get(r, w); err == nil {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return nil
		}
		return templates.Render(w, "login", nil)
	}
}

func LoginPost() Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		if _, err := session.Get(r, w); err == nil {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return nil
		}

		r.ParseForm()
		password := r.FormValue("password")
		if !pkg.CheckPasswordHash(password, config.Get().App.AdminPassword) {
			return templates.Render(w, "login", dto.LoginTemplateDTO{Message: "wrong credentials"})
		}
		if err := session.Start(r, w); err != nil {
			return err
		}
		http.Redirect(w, r, "/admin", http.StatusFound)
		return nil
	}
}
