package endpoints

import (
	"microblog-app/internal/dto"
	"microblog-app/internal/templates"
	"net/http"
)

func Admin() Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		return templates.Render(w, "admin", dto.AdminTemplateDTO{})
	}
}
