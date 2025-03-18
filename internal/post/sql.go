package post

import (
	"microblog-app/pkg"
	"time"
)

type PostSQL struct {
	ID        uint64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToSQL(post Post) PostSQL {
	return PostSQL{
		ID:        post.id.Value(),
		Content:   post.content.Value(),
		CreatedAt: post.createdAt.Value(),
		UpdatedAt: post.updatedAt.Value(),
	}
}

func (p PostSQL) Domain() (Post, error) {
	problems := pkg.ValidationError{}
	newID, err := NewID(p.ID)
	if err != nil {
		problems["id"] = err.Error()
	}
	newContent, err := NewContent(p.Content)
	if err != nil {
		problems["content"] = err.Error()
	}
	newCreatedAt, err := NewCreatedAt(p.CreatedAt)
	if err != nil {
		problems["createdAt"] = err.Error()
	}
	newUpdatedAt, err := NewUpdatedAt(p.UpdatedAt)
	if err != nil {
		problems["updatedAt"] = err.Error()
	}
	if len(problems) != 0 {
		return Post{}, problems
	}
	return Post{
		id:        newID,
		content:   newContent,
		createdAt: newCreatedAt,
		updatedAt: newUpdatedAt,
	}, nil
}
