package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
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

func (ps *postService) Create(ctx context.Context, req request.CreatePostRequest) (*response.PostResponse, error) {
	post := &entity.Post{
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
		UserId:  1, // TODO: get from auth context
	}

	value, err := ps.postRepository.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	return &response.PostResponse{
		ID:        value.ID,
		Content:   value.Content,
		Title:     value.Title,
		UserId:    value.UserId,
		Tags:      value.Tags,
		Version:   value.Version,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}, nil
}

func (ps *postService) GetById(ctx context.Context, postID int64) (*response.PostResponse, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		return nil, err
	}

	return &response.PostResponse{
		ID:        post.ID,
		Content:   post.Content,
		Title:     post.Title,
		UserId:    post.UserId,
		Tags:      post.Tags,
		Version:   post.Version,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (ps *postService) GetByIdWithComments(ctx context.Context, postID int64) (*response.PostResponse, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		return nil, err
	}

	comments, err := ps.commentRepository.GetByPostId(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	post.Comment = comments

	commentResponses := make([]response.CommentResponse, 0, len(comments))

	for _, cmt := range comments {
		commentResponses = append(commentResponses, response.CommentResponse{
			ID:      cmt.ID,
			PostId:  cmt.PostId,
			UserId:  cmt.UserId,
			Content: cmt.Content,
			User: response.UserResponse{
				ID:       cmt.User.ID,
				Username: cmt.User.Username,
			},
			CreatedAt: cmt.CreatedAt,
		})
	}

	return &response.PostResponse{
		ID:              post.ID,
		Content:         post.Content,
		Title:           post.Title,
		UserId:          post.UserId,
		Tags:            post.Tags,
		Version:         post.Version,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		CommentResponse: commentResponses,
	}, nil
}

func (ps *postService) Delete(ctx context.Context, postID int64) error {
	return ps.postRepository.Delete(ctx, postID)
}

func (ps *postService) Update(ctx context.Context, postID int64, req request.UpdatePostRequest) (*response.PostResponse, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &response.PostResponse{
		ID:        updatedPost.ID,
		Content:   updatedPost.Content,
		Title:     updatedPost.Title,
		UserId:    updatedPost.UserId,
		Tags:      updatedPost.Tags,
		Version:   updatedPost.Version,
		CreatedAt: updatedPost.CreatedAt,
		UpdatedAt: updatedPost.UpdatedAt,
	}, nil
}
