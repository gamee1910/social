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
	postRepo repository.PostRepository, commentRepo repository.CommentRepository, logger *logger.Logger,
) service.PostService {
	return &postService{
		postRepository:    postRepo,
		commentRepository: commentRepo,
		logger:            logger,
	}
}

func (postService *postService) Create(ctx context.Context, postRequest request.CreatePostRequest) (*response.PostResponse, error) {
	req := &entity.Post{
		Title:   postRequest.Title,
		Content: postRequest.Content,
		Tags:    postRequest.Tags,
		UserId:  postRequest.UserID, // lấy từ JWT claims qua handler
	}

	post, err := postService.postRepository.Create(ctx, req)
	if err != nil {
		postService.logger.Error("failed to create post", "error", err)
		return nil, err
	}

	postService.logger.Info("post created successfully", "postID", post.ID, "userID", post.UserId)

	return &response.PostResponse{
		ID:        post.ID,
		Content:   post.Content,
		Title:     post.Title,
		UserId:    post.UserId,
		Tags:      post.Tags,
		Version:   post.Version,
		CreatedAt: post.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		UpdatedAt: post.UpdatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
	}, nil
}

func (postService *postService) GetById(ctx context.Context, postID int64) (*response.PostResponse, error) {
	post, err := postService.postRepository.GetById(ctx, postID)
	if err != nil {
		postService.logger.Error("failed to get post by id", "postID", postID, "error", err)
		return nil, err
	}

	postService.logger.Info("post retrieved successfully", "postID", post.ID)

	return &response.PostResponse{
		ID:        post.ID,
		Content:   post.Content,
		Title:     post.Title,
		UserId:    post.UserId,
		Tags:      post.Tags,
		Version:   post.Version,
		CreatedAt: post.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		UpdatedAt: post.UpdatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
	}, nil
}

func (postService *postService) GetByIdWithComments(ctx context.Context, postID int64) (*response.PostResponse, error) {
	post, err := postService.postRepository.GetById(ctx, postID)
	if err != nil {
		postService.logger.Error("failed to get post by id", "postID", postID, "error", err)
		return nil, err
	}

	comments, err := postService.commentRepository.GetByPostId(ctx, post.ID)
	if err != nil {
		postService.logger.Error("failed to get comments by post id", "postID", post.ID, "error", err)
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
			CreatedAt: cmt.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		})
	}

	postService.logger.Info("post with comments retrieved successfully", "postID", post.ID, "commentCount", len(comments))

	return &response.PostResponse{
		ID:              post.ID,
		Content:         post.Content,
		Title:           post.Title,
		UserId:          post.UserId,
		Tags:            post.Tags,
		Version:         post.Version,
		CreatedAt:       post.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		UpdatedAt:       post.UpdatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		CommentResponse: commentResponses,
	}, nil
}

func (postService *postService) Delete(ctx context.Context, postID int64) error {
	if err := postService.postRepository.Delete(ctx, postID); err != nil {
		postService.logger.Error("failed to delete post", "postID", postID, "error", err)
		return err
	}

	postService.logger.Info("post deleted successfully", "postID", postID)
	return nil
}

func (postService *postService) Update(ctx context.Context, postID int64, req request.UpdatePostRequest) (*response.PostResponse, error) {
	post, err := postService.postRepository.GetById(ctx, postID)
	if err != nil {
		postService.logger.Error("failed to get post for update", "postID", postID, "error", err)
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

	updatedPost, err := postService.postRepository.Update(ctx, postID, post)
	if err != nil {
		postService.logger.Error("failed to update post", "postID", postID, "error", err)
		return nil, err
	}

	postService.logger.Info("post updated successfully", "postID", updatedPost.ID)

	return &response.PostResponse{
		ID:        updatedPost.ID,
		Content:   updatedPost.Content,
		Title:     updatedPost.Title,
		UserId:    updatedPost.UserId,
		Tags:      updatedPost.Tags,
		Version:   updatedPost.Version,
		CreatedAt: updatedPost.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		UpdatedAt: updatedPost.UpdatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
	}, nil
}

func (postService *postService) GetFeed(ctx context.Context, input service.GetFeedInput) ([]*response.PostWithMetaData, error) {
	if input.Limit <= 0 {
		input.Limit = 20
	}

	if input.Limit > 20 {
		input.Limit = 20
	}

	if input.Offset < 0 {
		input.Offset = 0
	}

	if input.Sort != "asc" && input.Sort != "desc" {
		input.Sort = "desc"
	}

	query := repository.FeedQuery{
		UserID: input.UserID,
		Limit:  input.Limit,
		Offset: input.Offset,
		Sort:   input.Sort,
		Search: input.Search,
		Tags:   input.Tags,
	}

	posts, err := postService.postRepository.GetFeed(ctx, query)
	if err != nil {
		return nil, err
	}

	feed := make([]*response.PostWithMetaData, 0, len(posts))

	for _, post := range posts {
		items := &response.PostWithMetaData{
			PostResponse: response.PostResponse{
				ID:        post.ID,
				Content:   post.Content,
				Title:     post.Title,
				UserId:    post.UserId,
				Tags:      post.Tags,
				Version:   post.Version,
				CreatedAt: post.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
				UpdatedAt: post.UpdatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
			},
			CommentsCount: post.CommentsCount,
		}
		feed = append(feed, items)
	}

	return feed, nil
}
