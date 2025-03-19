package main

import (
	"embed"
	"fmt"
	"microblog-app/internal/app"
	"microblog-app/pkg"
	"os"
)

//go:embed templates/*
var templates embed.FS

//go:embed assets/*
var assets embed.FS

//go:embed migrations/*
var migrations embed.FS

func main() {
	args := os.Args[1:]
	if len(args) == 2 {
		if args[0] == "bcrypt" {
			password, err := pkg.HashPassword(args[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			fmt.Fprintln(os.Stdout, password)
			return
		}
	}

	app := app.New(templates, assets, migrations)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
