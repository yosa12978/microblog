package endpoints

import (
	"bytes"
	"fmt"
	"microblog-app/internal/config"
	"microblog-app/internal/services"
	"net/http"

	"github.com/gorilla/feeds"
	"github.com/yuin/goldmark"
)

func GetFeed(postService services.PostService) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) error {
		conf := config.Get()
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		hostname := conf.Server.Hostname
		if hostname == "" {
			hostname = r.URL.Host
		}
		baseLink := fmt.Sprintf("%s://%s", scheme, hostname)
		feed := &feeds.Feed{
			Title:       conf.App.Name,
			Link:        &feeds.Link{Href: baseLink},
			Description: conf.App.Desc,
		}

		feed.Items = []*feeds.Item{}

		posts, err := postService.GetFeed(r.Context())
		for _, v := range posts {
			title := v.Content
			if len(v.Content) > 20 {
				title = v.Content[:20] + "..."
			}
			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(v.Content), &buf); err != nil {
				panic(err) // don't like this
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
			return err
		}
		w.Header().Set("Content-Type", "application/atom+xml")
		w.WriteHeader(200)
		w.Write([]byte(atomFeed))
		return nil
	}
}
