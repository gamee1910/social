package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/store"
)

type CommentService struct {
	commentStore store.CommentRepository
}

func NewCommentService(commentStore store.CommentRepository) *CommentService {
	return &CommentService{
		commentStore: commentStore,
	}
}

func (cs *CommentService) GetByPostId(ctx context.Context, postID int64) ([]entity.Comment, error) {
	return cs.commentStore.GetByPostId(ctx, postID)
}
