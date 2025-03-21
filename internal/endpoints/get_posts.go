package endpoints

import (
	"microblog-app/internal/config"
	"microblog-app/internal/dto"
	"microblog-app/internal/services"
	"microblog-app/internal/templates"
	"net/http"
	"strconv"
)

func GetPosts(postService services.PostService) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		pageStr := r.URL.Query().Get("page")
		sizeStr := r.URL.Query().Get("size")

		page, err := strconv.ParseUint(pageStr, 10, 32)
		if err != nil {
			page = 1
		}
		size, err := strconv.ParseUint(sizeStr, 10, 32)
		if err != nil {
			size = 20
		}

		posts, err := postService.GetAll(r.Context(), uint(page), uint(size))
		if err != nil {
			return err
		}

		payload := dto.BlogTemplateDTO{Posts: posts, Description: config.Get().App.Desc}
		return templates.Render(w, "blog", payload)
	}
}
