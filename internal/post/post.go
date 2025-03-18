package post

import (
	"microblog-app/pkg"
	"time"
)

type Post struct {
	id        ID
	content   Content
	createdAt CreatedAt
	updatedAt UpdatedAt
}

func New(id ID, content Content, createdAt CreatedAt, updatedAt UpdatedAt) Post {
	return Post{
		id:        id,
		content:   content,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func NewFromPrimitives(id uint64, content string) (Post, error) {
	problems := pkg.ValidationError{}
	newID, err := NewID(id)
	if err != nil {
		problems["id"] = err.Error()
	}
	newContent, err := NewContent(content)
	if err != nil {
		problems["content"] = err.Error()
	}
	now := time.Now().UTC()
	newCreatedAt, _ := NewCreatedAt(now)
	newUpdatedAt, _ := NewUpdatedAt(now)
	if len(problems) != 0 {
		return Post{}, problems
	}
	return New(newID, newContent, newCreatedAt, newUpdatedAt), nil
}

func (p *Post) ID() ID {
	return p.id
}

func (p *Post) Content() Content {
	return p.content
}

func (p *Post) CreatedAt() CreatedAt {
	return p.createdAt
}

func (p *Post) UpdatedAt() UpdatedAt {
	return p.updatedAt
}

func (p *Post) SetID(id ID) error {
	p.id = id
	return nil
}

func (p *Post) SetContent(content Content) error {
	p.content = content
	return nil
}

func (p *Post) SetCreatedAt(createdAt CreatedAt) error {
	p.createdAt = createdAt
	return nil
}

func (p *Post) SetUpdatedAt(updatedAt UpdatedAt) error {
	p.updatedAt = updatedAt
	return nil
}
