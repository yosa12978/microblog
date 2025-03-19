package post

import (
	"context"
	"microblog-app/pkg"
)

type Repo interface {
	GetAll(ctx context.Context, page, size uint) (pkg.Page[Post], error)
	GetFeed(ctx context.Context) ([]Post, error)
	GetByID(ctx context.Context, id ID) (Post, error)
	Create(ctx context.Context, p Post) (ID, error)
	Update(ctx context.Context, p Post) (ID, error)
	Delete(ctx context.Context, id ID) (ID, error)
}
