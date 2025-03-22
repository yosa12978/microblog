package endpoints

import (
	"microblog-app/internal/dto"
	"microblog-app/internal/services"
	"microblog-app/internal/templates"
	"net/http"
	"strconv"
)

func DeletePost(postService services.PostService) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseForm()
		idStr := r.FormValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return templates.Render(w, "admin",
				dto.AdminTemplateDTO{DeleteErrorMessage: "id: invalid id provided"})
		}
		if _, err := postService.Delete(r.Context(), id); err != nil {
			return templates.Render(w, "admin",
				dto.AdminTemplateDTO{DeleteErrorMessage: err.Error()})
		}
		return templates.Render(w, "admin",
			dto.AdminTemplateDTO{DeleteSuccessMessage: "success"})
	}
}
