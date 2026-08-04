package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
)

type commentService struct {
	commentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) service.CommentService {
	return &commentService{commentRepository: commentRepository}
}

func (cs *commentService) GetByPostId(ctx context.Context, postID int64) ([]entity.Comment, error) {
	return cs.commentRepository.GetByPostId(ctx, postID)
}
