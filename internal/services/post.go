package services

import (
	"context"
	"log/slog"
	"microblog-app/internal/dto"
	"microblog-app/internal/post"
	"microblog-app/pkg"
)

type PostService interface {
	GetAll(ctx context.Context, page, size uint) (*pkg.Page[dto.PostResponse], error)
	GetFeed(ctx context.Context) ([]dto.PostResponse, error)
	GetByID(ctx context.Context, id uint64) (*dto.PostResponse, error)

	Create(ctx context.Context, req dto.PostCreateRequest) (uint64, error)
	Update(ctx context.Context, id uint64, req dto.PostUpdateRequest) (uint64, error)
	Delete(ctx context.Context, id uint64) (uint64, error)
}

type postService struct {
	repo   post.Repo
	logger *slog.Logger
}

func NewPostService(repo post.Repo, logger *slog.Logger) PostService {
	return &postService{
		repo:   repo,
		logger: logger,
	}
}

func (p *postService) Create(ctx context.Context, req dto.PostCreateRequest) (uint64, error) {
	post, err := req.Domain()
	if err != nil {
		return 0, err
	}
	id, err := p.repo.Create(ctx, post)
	return id.Value(), err
}

func (p *postService) Delete(ctx context.Context, id uint64) (uint64, error) {
	domainID, err := post.NewID(id)
	if err != nil {
		return 0, err
	}
	respID, err := p.repo.Delete(ctx, domainID)
	return respID.Value(), err
}

func (p *postService) GetAll(
	ctx context.Context,
	page uint,
	size uint,
) (*pkg.Page[dto.PostResponse], error) {
	posts, err := p.repo.GetAll(ctx, page, size)
	if err != nil {
		return nil, err
	}
	result := pkg.NewPage([]dto.PostResponse{}, posts.Total, posts.Current, posts.Size)
	for _, v := range posts.Content {
		result.Content = append(result.Content, dto.NewPostResponse(v))
	}
	return &result, nil
}

func (p *postService) GetFeed(ctx context.Context) ([]dto.PostResponse, error) {
	posts, err := p.repo.GetFeed(ctx)
	if err != nil {
		return []dto.PostResponse{}, err
	}
	res := []dto.PostResponse{}
	for _, v := range posts {
		res = append(res, dto.NewPostResponse(v))
	}
	return res, nil
}

func (p *postService) GetByID(ctx context.Context, id uint64) (*dto.PostResponse, error) {
	post, err := p.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := dto.NewPostResponse(post)
	return &response, nil
}

func (p *postService) getByID(ctx context.Context, id uint64) (post.Post, error) {
	domainID, err := post.NewID(id)
	if err != nil {
		return post.Post{}, err
	}
	return p.repo.GetByID(ctx, domainID)
}

func (p *postService) Update(
	ctx context.Context,
	id uint64,
	req dto.PostUpdateRequest,
) (uint64, error) {
	post, err := p.getByID(ctx, id)
	if err != nil {
		return 0, err
	}
	respID, err := p.repo.Update(ctx, post)
	return respID.Value(), err
}
