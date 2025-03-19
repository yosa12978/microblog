package router

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"microblog-app/internal/config"
	"microblog-app/internal/dto"
	"microblog-app/internal/middleware"
	"microblog-app/internal/services"
	"microblog-app/internal/session"
	"microblog-app/internal/templates"
	"microblog-app/pkg"
	"net/http"
	"strconv"

	"github.com/gorilla/feeds"
	"github.com/yuin/goldmark"
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
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	router.Handle("/assets/", http.FileServerFS(opts.AssetsFS))

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		sizeStr := r.URL.Query().Get("size")

		page, err := strconv.ParseUint(pageStr, 10, 32)
		if err != nil {
			page = 1
		}
		size, err := strconv.ParseUint(sizeStr, 10, 32)
		if err != nil {
			size = 20
		}

		posts, err := opts.PostService.GetAll(r.Context(), uint(page), uint(size))
		payload := dto.BlogTemplateDTO{Posts: posts, Description: opts.Config.App.Desc}
		templates.Render(w, "blog", payload)
	})

	router.HandleFunc("GET /feed", func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		hostname := opts.Config.Server.Hostname
		if hostname == "" {
			hostname = r.URL.Host
		}
		baseLink := fmt.Sprintf("%s://%s", scheme, hostname)
		feed := &feeds.Feed{
			Title:       opts.Config.App.Name,
			Link:        &feeds.Link{Href: baseLink},
			Description: opts.Config.App.Desc,
		}

		feed.Items = []*feeds.Item{}

		posts, err := opts.PostService.GetFeed(r.Context())
		for _, v := range posts {
			title := v.Content
			if len(v.Content) > 20 {
				title = v.Content[:20] + "..."
			}
			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(v.Content), &buf); err != nil {
				panic(err)
			}
			content := buf.String()
			item := &feeds.Item{
				Id:    fmt.Sprintf("%v", v.ID),
				Title: title,
				Link: &feeds.Link{
					Href: fmt.Sprintf("%s/posts/%v", baseLink, v.ID),
				},
				Created: v.CreatedAt,
				Updated: v.UpdatedAt,
				Content: content,
			}
			feed.Items = append(feed.Items, item)
		}

		atomFeed, err := feed.ToAtom()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/atom+xml")
		w.WriteHeader(200)
		w.Write([]byte(atomFeed))
	})

	router.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		post, err := opts.PostService.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		templates.Render(w, "post", dto.PostTemplateDTO{Post: *post})
	})

	router.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		if _, err := session.Get(r, w); err == nil {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}

		templates.Render(w, "login", nil)
	})

	router.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		if _, err := session.Get(r, w); err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if err := session.End(r, w); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	router.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		if _, err := session.Get(r, w); err == nil {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}

		r.ParseForm()
		password := r.FormValue("password")
		if !pkg.CheckPasswordHash(password, opts.Config.App.AdminPassword) {
			templates.Render(w, "login", dto.LoginTemplateDTO{Message: "wrong credentials"})
			return
		}
		if err := session.Start(r, w); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusFound)
	})

	router.Handle(
		"GET /admin",
		middleware.Admin()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("admin page"))
		})),
	)
}
