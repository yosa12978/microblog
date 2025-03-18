package main

import (
	"embed"
	"microblog-app/internal/app"
)

//go:embed templates/*
var templates embed.FS

//go:embed assets/*
var assets embed.FS

//go:embed migrations/*
var migrations embed.FS

func main() {
	app := app.New(templates, assets, migrations)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
