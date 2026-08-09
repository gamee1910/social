package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
	"github.com/gamee1910/social/pkg/logger"
)

type commentService struct {
	commentRepository repository.CommentRepository
	logger            *logger.Logger
}

func NewCommentService(
	commentRepository repository.CommentRepository,
	logger *logger.Logger,
) service.CommentService {
	return &commentService{
		commentRepository: commentRepository,
		logger:            logger,
	}
}

func (cs *commentService) GetByPostId(ctx context.Context, postID int64) ([]response.CommentResponse, error) {
	entities, err := cs.commentRepository.GetByPostId(ctx, postID)
	if err != nil {
		cs.logger.Error("failed to get comments by post id", "postID", postID, "error", err)
		return nil, err
	}

	responses := make([]response.CommentResponse, 0, len(entities))

	for _, cmt := range entities {
		responses = append(responses, response.CommentResponse{
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

	cs.logger.Info("comments retrieved successfully", "postID", postID, "count", len(responses))

	return responses, nil
}
