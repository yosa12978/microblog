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

type LoginTemplateDTO struct {
	Message string
}

type AdminTemplateDTO struct {
	CreateSuccessMessage string
	CreateErrorMessage   string
	UpdateSuccessMessage string
	UpdateErrorMessage   string
	DeleteSuccessMessage string
	DeleteErrorMessage   string
	PinSuccessMessage    string
	PinErrorMessage      string
}
