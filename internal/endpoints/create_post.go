package endpoints

import (
	"microblog-app/internal/dto"
	"microblog-app/internal/services"
	"microblog-app/internal/templates"
	"net/http"
)

func CreatePost(postService services.PostService) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseForm()
		content := r.FormValue("content")
		pin := r.FormValue("pin")
		if _, err := postService.Create(r.Context(),
			dto.PostCreateRequest{Content: content, Pinned: pin == "on"}); err != nil {
			return templates.Render(w, "admin",
				dto.AdminTemplateDTO{CreateErrorMessage: err.Error()})
		}
		return templates.Render(w, "admin",
			dto.AdminTemplateDTO{CreateSuccessMessage: "success"})
	}
}
