package dto

import (
	"microblog-app/pkg"
)

type BlogTemplateDTO struct {
	Posts       *pkg.Page[PostResponse]
	Description string
}

type PostTemplateDTO struct {
	Post PostResponse
}
