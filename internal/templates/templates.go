package templates

import (
	"html/template"
	"io"
	"io/fs"
	"microblog-app/internal/config"
)

var (
	templatesFS fs.FS
)

func Init(templates fs.FS) {
	templatesFS = templates
}

func Render(w io.Writer, name string, data any) error {
	if templatesFS == nil {
		panic("templateFS is nil")
	}
	template, err := template.ParseFS(
		templatesFS,
		"templates/pages/"+name+".tmpl",
		"templates/top.tmpl",
		"templates/bottom.tmpl",
	)
	if err != nil {
		return err
	}
	conf := config.Get()
	payload := map[string]any{
		"Title":   conf.App.Name,
		"Bottom":  conf.App.Bottom,
		"Content": data,
	}
	return template.Execute(w, payload)
}
