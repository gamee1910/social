package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
	"github.com/gamee1910/social/pkg/logger"
)

type postService struct {
	postRepository    repository.PostRepository
	commentRepository repository.CommentRepository
	logger            *logger.Logger
}

func NewPostService(
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	logger *logger.Logger,
) service.PostService {
	return &postService{
		postRepository:    postRepo,
		commentRepository: commentRepo,
		logger:            logger,
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
		ps.logger.Error("failed to create post", "error", err)
		return nil, err
	}

	ps.logger.Info("post created successfully", "postID", value.ID, "userID", value.UserId)

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
		ps.logger.Error("failed to get post by id", "postID", postID, "error", err)
		return nil, err
	}

	ps.logger.Info("post retrieved successfully", "postID", post.ID)

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
		ps.logger.Error("failed to get post by id", "postID", postID, "error", err)
		return nil, err
	}

	comments, err := ps.commentRepository.GetByPostId(ctx, post.ID)
	if err != nil {
		ps.logger.Error("failed to get comments by post id", "postID", post.ID, "error", err)
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

	ps.logger.Info("post with comments retrieved successfully", "postID", post.ID, "commentCount", len(comments))

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
	if err := ps.postRepository.Delete(ctx, postID); err != nil {
		ps.logger.Error("failed to delete post", "postID", postID, "error", err)
		return err
	}

	ps.logger.Info("post deleted successfully", "postID", postID)
	return nil
}

func (ps *postService) Update(ctx context.Context, postID int64, req request.UpdatePostRequest) (*response.PostResponse, error) {
	post, err := ps.postRepository.GetById(ctx, postID)
	if err != nil {
		ps.logger.Error("failed to get post for update", "postID", postID, "error", err)
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
		ps.logger.Error("failed to update post", "postID", postID, "error", err)
		return nil, err
	}

	ps.logger.Info("post updated successfully", "postID", updatedPost.ID)

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
