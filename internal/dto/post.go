package dto

import (
	"microblog-app/internal/post"
	"microblog-app/pkg"
	"time"
)

type PostResponse struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	Pinned    bool      `json:"pinned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPostResponse(domain post.Post) PostResponse {
	return PostResponse{
		ID:        domain.ID().Value(),
		Content:   domain.Content().Value(),
		Pinned:    domain.Pinned().Value(),
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
	Pinned  bool   `json:"pinned"`
}

func (req PostUpdateRequest) Apply(model *post.Post) error {
	problems := pkg.ValidationError{}
	newContent, err := post.NewContent(req.Content)
	if err != nil {
		problems["content"] = err.Error()
	}
	newPinned, err := post.NewPinned(req.Pinned)
	if err != nil {
		problems["pinned"] = err.Error()
	}
	if len(problems) != 0 {
		return problems
	}
	model.SetPinned(newPinned)
	model.SetContent(newContent)
	return nil
}
