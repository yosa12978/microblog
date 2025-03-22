package endpoints

import (
	"microblog-app/internal/dto"
	"microblog-app/internal/services"
	"microblog-app/internal/templates"
	"net/http"
	"strconv"
)

func PinPost(postService services.PostService) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseForm()
		idStr := r.FormValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return templates.Render(w, "admin",
				dto.AdminTemplateDTO{PinErrorMessage: "id: invalid id provided"})
		}
		if err := postService.Pin(r.Context(), id); err != nil {
			return templates.Render(w, "admin",
				dto.AdminTemplateDTO{PinErrorMessage: err.Error()})
		}
		return templates.Render(w, "admin",
			dto.AdminTemplateDTO{PinSuccessMessage: "success"})
	}
}
