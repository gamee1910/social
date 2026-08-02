package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/store"
)

type CommentService struct {
	commentRepository store.CommentRepository
}

func NewCommentService(commentStore store.CommentRepository) *CommentService {
	return &CommentService{
		commentRepository: commentStore,
	}
}

func (cs *CommentService) GetByPostId(ctx context.Context, postID int64) ([]entity.Comment, error) {
	return cs.commentRepository.GetByPostId(ctx, postID)
}
