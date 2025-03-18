package app

import (
	"io/fs"
	"log/slog"
	"microblog-app/internal/config"
	"microblog-app/internal/data"
	"microblog-app/internal/repos"
	"microblog-app/internal/router"
	"microblog-app/internal/services"
	"net/http"
)

func newServer(addr string, logger *slog.Logger, assets fs.FS, config config.Config) http.Server {
	postRepo := repos.NewPostRepoSQL(data.Postgres(), logger)
	postService := services.NewPostService(postRepo, logger)
	router := router.New(
		&router.RouterOptions{
			Logger:      logger,
			AssetsFS:    assets,
			Config:      config,
			PostService: postService,
		},
	)
	return http.Server{
		Handler: router,
		Addr:    addr,
	}
}
