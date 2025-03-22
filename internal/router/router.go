package router

import (
	"io/fs"
	"log/slog"
	"microblog-app/internal/config"
	"microblog-app/internal/endpoints"
	"microblog-app/internal/middleware"
	"microblog-app/internal/services"
	"net/http"
)

type RouterOptions struct {
	Logger      *slog.Logger
	AssetsFS    fs.FS
	Config      config.Config
	PostService services.PostService
}

func New(opts *RouterOptions) http.Handler {
	router := http.NewServeMux()
	addRoutes(router, opts)
	handler := middleware.Pipeline(
		router,
		middleware.Logger(opts.Logger),
		middleware.StripSlash(),
		middleware.Recovery(opts.Logger),
	)
	return handler
}

func addRoutes(router *http.ServeMux, opts *RouterOptions) {
	router.HandleFunc("/health", endpoints.Health())
	router.Handle("/assets/", http.FileServerFS(opts.AssetsFS))
	router.HandleFunc("GET /{$}", endpoints.GetPosts(opts.PostService).Unwrap())
	router.HandleFunc("GET /feed", endpoints.GetFeed(opts.PostService).Unwrap())
	router.HandleFunc("GET /posts/{id}", endpoints.GetPostByID(opts.PostService).Unwrap())

	router.HandleFunc("GET /login", endpoints.LoginGet().Unwrap())
	router.HandleFunc("GET /logout", endpoints.Logout().Unwrap())
	router.HandleFunc("POST /login", endpoints.LoginPost().Unwrap())

	addAdminRoutes(router, opts)
}

func addAdminRoutes(router *http.ServeMux, opts *RouterOptions) {
	router.Handle(
		"GET /admin",
		middleware.Admin()(endpoints.Admin().Unwrap()),
	)
	router.Handle(
		"POST /admin/create-post",
		middleware.Admin()(endpoints.CreatePost(opts.PostService).Unwrap()),
	)
	router.Handle(
		"POST /admin/pin-post",
		middleware.Admin()(endpoints.PinPost(opts.PostService).Unwrap()),
	)
	router.Handle(
		"POST /admin/delete-post",
		middleware.Admin()(endpoints.DeletePost(opts.PostService).Unwrap()),
	)
	router.Handle(
		"POST /admin/update-post",
		middleware.Admin()(endpoints.UpdatePost(opts.PostService).Unwrap()),
	)
}
