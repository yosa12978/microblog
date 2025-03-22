package repos

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"microblog-app/internal/post"
	"microblog-app/pkg"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepoPGX struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostRepoPGX(db *pgxpool.Pool, logger *slog.Logger) post.Repo {
	return &postRepoPGX{
		db:     db,
		logger: logger,
	}
}

const (
	getPostsCountSQL = `
        SELECT COUNT(*) FROM posts;
    `
	getAllPostsSQL = `
        SELECT id, content, pinned, created_at, updated_at 
        FROM posts ORDER BY pinned DESC, created_at DESC LIMIT $1 OFFSET $2;
    `
	getPostsFeedSQL = `
        SELECT id, content, pinned, created_at, updated_at 
        FROM posts ORDER BY created_at DESC LIMIT 30;
    `
	getPostByIDSQL = `
        SELECT id, content, pinned, created_at, updated_at
        FROM posts WHERE id=$1;
    `
	createPostSQL = `
        INSERT INTO posts(content, created_at, updated_at, pinned)
        VALUES ($1, $2, $3, $4) RETURNING id;
    `
	updatePostSQL = `
        UPDATE posts SET content=$1, pinned=$3, updated_at=now()
        WHERE id=$2 RETURNING id;
    `
	deletePostSQL = `
        DELETE FROM posts WHERE id=$1;
    `
)

func (repo *postRepoPGX) GetAll(
	ctx context.Context,
	page uint,
	size uint,
) (pkg.Page[post.Post], error) {
	var count uint
	err := repo.db.QueryRow(ctx, getPostsCountSQL).Scan(&count)
	if err != nil {
		return pkg.Page[post.Post]{}, err
	}
	totalPages := uint(math.Ceil(float64(count) / float64(size)))
	rows, err := repo.db.Query(ctx, getAllPostsSQL, size, (page-1)*size)
	if err != nil {
		return pkg.Page[post.Post]{}, nil
	}

	var posts []post.Post
	for rows.Next() {
		var postSQL post.PostSQL
		if err := rows.Scan(
			&postSQL.ID,
			&postSQL.Content,
			&postSQL.Pinned,
			&postSQL.CreatedAt,
			&postSQL.UpdatedAt,
		); err != nil {
			repo.logger.Error("failed to scan post", "error", err.Error())
			// replace with corrupted post or smth
			continue
		}
		domain, err := postSQL.Domain()
		if err != nil {
			repo.logger.Error("failed to convert to domain", "error", err.Error())
			// replace with corrupted post or smth
			continue
		}
		posts = append(posts, domain)
	}

	return pkg.NewPage(posts, totalPages, page, size), nil
}

func (repo *postRepoPGX) GetFeed(ctx context.Context) ([]post.Post, error) {
	rows, err := repo.db.Query(ctx, getPostsFeedSQL)
	if err != nil {
		return []post.Post{}, nil
	}

	var posts []post.Post
	for rows.Next() {
		var postSQL post.PostSQL
		if err := rows.Scan(
			&postSQL.ID,
			&postSQL.Content,
			&postSQL.Pinned,
			&postSQL.CreatedAt,
			&postSQL.UpdatedAt,
		); err != nil {
			repo.logger.Error("failed to scan post", "error", err.Error())
			continue
		}
		domain, err := postSQL.Domain()
		if err != nil {
			repo.logger.Error("failed to convert to domain", "error", err.Error())
			continue
		}
		posts = append(posts, domain)
	}
	return posts, nil
}

func (repo *postRepoPGX) GetByID(ctx context.Context, id post.ID) (post.Post, error) {
	var postSQL post.PostSQL

	err := repo.db.QueryRow(ctx, getPostByIDSQL, id.Value()).
		Scan(&postSQL.ID, &postSQL.Content, &postSQL.Pinned, &postSQL.CreatedAt, &postSQL.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return post.Post{}, errors.New("post not found")
		}
		return post.Post{}, err
	}
	return postSQL.Domain()
}

func (repo *postRepoPGX) Create(ctx context.Context, p post.Post) (post.ID, error) {
	var lastInsertedID uint64
	err := repo.db.QueryRow(
		ctx,
		createPostSQL,
		p.Content().Value(),
		p.UpdatedAt().Value(),
		p.CreatedAt().Value(),
		p.Pinned().Value(),
	).Scan(&lastInsertedID)
	if err != nil {
		return 0, err
	}
	return post.NewID(lastInsertedID)
}

func (repo *postRepoPGX) Update(ctx context.Context, p post.Post) (post.ID, error) {
	var lastInsertedID uint64
	err := repo.db.QueryRow(
		ctx,
		updatePostSQL,
		p.Content().Value(),
		p.ID().Value(),
		p.Pinned().Value(),
	).Scan(&lastInsertedID)
	if err != nil {
		return 0, err
	}
	return post.NewID(lastInsertedID)
}

func (repo *postRepoPGX) Delete(ctx context.Context, id post.ID) (post.ID, error) {
	if _, err := repo.GetByID(ctx, id); err != nil {
		return 0, err
	}
	_, err := repo.db.Exec(
		ctx,
		deletePostSQL,
		id.Value(),
	)
	return id, err
}

func (repo *postRepoPGX) Pin(ctx context.Context, id post.ID) (post.ID, error) {
	post, err := repo.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}
	if err := post.SetPinned(!post.Pinned()); err != nil {
		return 0, err
	}
	return repo.Update(ctx, post)
}
