package repos

import (
	"context"
	"database/sql"
	"log/slog"
	"math"
	"microblog-app/internal/post"
	"microblog-app/pkg"
)

type postRepoSQL struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostRepoSQL(db *sql.DB, logger *slog.Logger) post.Repo {
	return &postRepoSQL{
		db:     db,
		logger: logger,
	}
}

const (
	getPostsCountSQL = `
        SELECT COUNT(*) FROM posts;
    `
	getAllPostsSQL = `
        SELECT id, content, created_at, updated_at 
        FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2;
    `
	getPostByIDSQL = `
        SELECT id, content, created_at, updated_at
        FROM posts WHERE id=$1;
    `
	createPostSQL = `
        INSERT INTO posts(id, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4);
    `
	updatePostSQL = `
        UPDATE posts SET content=$1, updated_at=now()
        WHERE id=$2;
    `
	deletePostSQL = `
        DELETE FROM posts WHERE id=$1;
    `
)

func (repo *postRepoSQL) GetAll(
	ctx context.Context,
	page uint,
	size uint,
) (pkg.Page[post.Post], error) {
	var count uint
	err := repo.db.QueryRowContext(ctx, getPostsCountSQL).Scan(&count)
	if err != nil {
		return pkg.Page[post.Post]{}, err
	}
	totalPages := uint(math.Ceil(float64(count / size)))
	rows, err := repo.db.QueryContext(ctx, getAllPostsSQL, size, (page-1)*size)
	if err != nil {
		return pkg.Page[post.Post]{}, nil
	}

	var posts []post.Post
	for rows.Next() {
		var postSQL post.PostSQL
		if err := rows.Scan(
			&postSQL.ID,
			&postSQL.Content,
			&postSQL.CreatedAt,
			&postSQL.UpdatedAt,
		); err != nil {
			return pkg.Page[post.Post]{}, err
		}
		domain, err := postSQL.Domain()
		if err != nil {
			return pkg.Page[post.Post]{}, err
		}
		posts = append(posts, domain)
	}

	return pkg.NewPage(posts, totalPages, page, size), nil
}

func (repo *postRepoSQL) GetByID(ctx context.Context, id post.ID) (post.Post, error) {
	var postSQL post.PostSQL

	err := repo.db.QueryRowContext(ctx, getPostByIDSQL, id.Value()).
		Scan(&postSQL.ID, &postSQL.Content, &postSQL.CreatedAt, &postSQL.UpdatedAt)
	if err != nil {
		return post.Post{}, err
	}
	return postSQL.Domain()
}

func (repo *postRepoSQL) Create(ctx context.Context, p post.Post) (post.ID, error) {
	result, err := repo.db.ExecContext(
		ctx,
		createPostSQL,
		p.ID().Value(),
		p.Content().Value(),
		p.UpdatedAt().Value(),
		p.CreatedAt().Value(),
	)
	if err != nil {
		return post.ID(0), err
	}
	idInt, err := result.LastInsertId()
	if err != nil {
		return post.ID(0), err
	}
	return post.NewID(uint64(idInt))
}

func (repo *postRepoSQL) Update(ctx context.Context, p post.Post) (post.ID, error) {
	result, err := repo.db.ExecContext(
		ctx,
		updatePostSQL,
		p.Content().Value(),
		p.ID().Value(),
	)
	if err != nil {
		return post.ID(0), err
	}
	idInt, err := result.LastInsertId()
	if err != nil {
		return post.ID(0), err
	}
	return post.NewID(uint64(idInt))
}

func (repo *postRepoSQL) Delete(ctx context.Context, id post.ID) (post.ID, error) {
	_, err := repo.db.ExecContext(
		ctx,
		deletePostSQL,
		id.Value(),
	)
	return id, err
}
