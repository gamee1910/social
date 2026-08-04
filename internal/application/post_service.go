package application

import (
	"context"
	"errors"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
)

type postService struct {
	postRepository    repository.PostRepository
	commentRepository repository.CommentRepository
}

func NewPostService(postRepo repository.PostRepository, commentRepo repository.CommentRepository) service.PostService {
	return &postService{
		postRepository:    postRepo,
		commentRepository: commentRepo,
	}
}

func (ps *postService) Create(ctx context.Context, req request.CreatePostRequest) (*entity.Post, error) {
	post := &entity.Post{
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
		UserId:  1, // TODO: get from auth context
	}

	if err := ps.postRepository.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *postService) GetById(ctx context.Context, postID int64) (*entity.Post, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		return nil, translatePostError(err)
	}
	return post, nil
}

func (ps *postService) GetByIdWithComments(ctx context.Context, postID int64) (*entity.Post, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		return nil, translatePostError(err)
	}

	comments, err := ps.commentRepository.GetByPostId(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	post.Comment = comments
	return post, nil
}

func (ps *postService) Delete(ctx context.Context, postID int64) error {
	if err := ps.postRepository.Delete(ctx, postID); err != nil {
		return translatePostError(err)
	}
	return nil
}

func (ps *postService) Update(ctx context.Context, postID int64, req request.UpdatePostRequest) (*entity.Post, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		return nil, translatePostError(err)
	}

	if req.Title != nil {
		post.Title = *req.Title
	}

	if req.Content != nil {
		post.Content = *req.Content
	}

	if len(req.Tags) > 0 {
		post.Tags = req.Tags
	}

	updatedPost, err := ps.postRepository.Update(ctx, postID, post)
	if err != nil {
		return nil, translatePostError(err)
	}

	return updatedPost, nil
}

func translatePostError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return &domain.NotFoundError{Resource: "post"}
	case errors.Is(err, domain.ErrVersionConflict):
		return &domain.ConflictError{Resource: "post"}
	default:
		return err
	}
}
