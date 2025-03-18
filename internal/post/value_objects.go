package post

import (
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type ID uint64
type Content string
type CreatedAt time.Time
type UpdatedAt time.Time

func NewID(id uint64) (ID, error) {
	return ID(id), nil
}

func (id ID) Value() uint64 {
	return uint64(id)
}

func NewContent(content string) (Content, error) {
	p := bluemonday.UGCPolicy()
	sanitized := p.Sanitize(content)
	trimmed := strings.TrimSpace(sanitized)
	if err := validateContent(trimmed); err != nil {
		return Content(""), err
	}
	return Content(trimmed), nil
}

func (content Content) Value() string {
	return string(content)
}

func NewCreatedAt(t time.Time) (CreatedAt, error) {
	return CreatedAt(t), nil
}

func (createdAt CreatedAt) Value() time.Time {
	return time.Time(createdAt)
}

func NewUpdatedAt(t time.Time) (UpdatedAt, error) {
	return UpdatedAt(t), nil
}

func (updatedAt UpdatedAt) Value() time.Time {
	return time.Time(updatedAt)
}
