package dto

import (
	"microblog-app/internal/post"
	"time"
)

type PostResponse struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPostResponse(domain post.Post) PostResponse {
	return PostResponse{
		ID:        domain.ID().Value(),
		Content:   domain.Content().Value(),
		CreatedAt: domain.CreatedAt().Value(),
		UpdatedAt: domain.UpdatedAt().Value(),
	}
}

type PostCreateRequest struct {
	Content string `json:"content"`
}

func (req PostCreateRequest) Domain() (post.Post, error) {
	return post.NewFromPrimitives(0, req.Content)
}

type PostUpdateRequest struct {
	Content string `json:"content"`
}

func (req PostUpdateRequest) Apply(model *post.Post) error {
	if req.Content != "" {
		newContent, err := post.NewContent(req.Content)
		if err != nil {
			return err
		}
		model.SetContent(newContent)
	}
	return nil
}
