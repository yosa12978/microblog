package endpoints

import (
	"microblog-app/internal/dto"
	"microblog-app/internal/services"
	"microblog-app/internal/templates"
	"net/http"
	"strconv"
)

func GetPostByID(postService services.PostService) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(400)
			return err
		}
		post, err := postService.GetByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(404)
			return err
		}
		return templates.Render(w, "post", dto.PostTemplateDTO{Post: *post})
	}
}
