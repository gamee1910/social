package service

import "github.com/gamee1910/social/internal/store"

type Service struct {
	PostsService   *PostService
	CommentService *CommentService
}

func NewService(storage *store.Storage) *Service {
	return &Service{
		PostsService:   NewPostService(storage.Posts, storage.Comments),
		CommentService: NewCommentService(storage.Comments),
	}
}
