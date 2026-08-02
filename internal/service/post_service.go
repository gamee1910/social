package service

import (
	"context"
	"errors"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/dto"
	"github.com/gamee1910/social/internal/store"
)

type PostService struct {
	postRepository    store.PostRepository
	commentRepository store.CommentRepository
}

func NewPostService(postStore store.PostRepository, commentStore store.CommentRepository) *PostService {
	return &PostService{
		postRepository:    postStore,
		commentRepository: commentStore,
	}
}

func (ps *PostService) Create(ctx context.Context, req dto.CreatePostRequest) (*entity.Post, error) {
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

func (ps *PostService) GetById(ctx context.Context, postID int64) (*entity.Post, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		return nil, translatePostError(err)
	}
	return post, nil
}

func (ps *PostService) GetByIdWithComments(ctx context.Context, postID int64) (*entity.Post, error) {
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

func (ps *PostService) Delete(ctx context.Context, postID int64) error {
	if err := ps.postRepository.Delete(ctx, postID); err != nil {
		return translatePostError(err)
	}
	return nil
}

func (ps *PostService) Update(ctx context.Context, postID int64, req dto.UpdatePostRequest) (*entity.Post, error) {
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
	case errors.Is(err, store.ErrNotFound):
		return &domain.NotFoundError{Resource: "post"}
	case errors.Is(err, store.ErrVersionConflict):
		return &domain.ConflictError{Resource: "post"}
	default:
		return err
	}
}
