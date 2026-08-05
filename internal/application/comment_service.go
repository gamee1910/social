package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
)

type commentService struct {
	commentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) service.CommentService {
	return &commentService{commentRepository: commentRepository}
}

func (cs *commentService) GetByPostId(ctx context.Context, postID int64) ([]response.CommentResponse, error) {
	entities, err := cs.commentRepository.GetByPostId(ctx, postID)
	if err != nil {
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

	return responses, nil
}
