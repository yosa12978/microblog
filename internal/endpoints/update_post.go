package endpoints

import (
	"microblog-app/internal/dto"
	"microblog-app/internal/services"
	"microblog-app/internal/templates"
	"net/http"
	"strconv"
)

func UpdatePost(postService services.PostService) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseForm()
		content := r.FormValue("content")
		pin := r.FormValue("pin")
		idStr := r.FormValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return templates.Render(w, "admin",
				dto.AdminTemplateDTO{UpdateErrorMessage: err.Error()})
		}
		_, err = postService.Update(r.Context(), id,
			dto.PostUpdateRequest{Content: content, Pinned: pin == "on"})
		if err != nil {
			return templates.Render(w, "admin",
				dto.AdminTemplateDTO{UpdateErrorMessage: err.Error()})
		}
		return templates.Render(w, "admin",
			dto.AdminTemplateDTO{UpdateSuccessMessage: "success"})
	}
}
