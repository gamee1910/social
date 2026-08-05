package di

import (
	"database/sql"

	"github.com/gamee1910/social/internal/application"
	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/infrastructure/persistences/postgres"
	"github.com/gamee1910/social/internal/interfaces/http/handler"
	"github.com/gamee1910/social/pkg/logger"
)

type Container struct {
	cfg    *config.Configuration
	db     *sql.DB
	logger *logger.Logger

	//Handlers
	postHandler    *handler.PostHandler
	userHandler    *handler.UserHandler
	commentHandler *handler.CommentHandler
}

func (c *Container) PostHandler() *handler.PostHandler {
	return c.postHandler
}

func (c *Container) UserHandler() *handler.UserHandler {
	return c.userHandler
}

func (c *Container) CommentHandler() *handler.CommentHandler {
	return c.commentHandler
}

func NewContainer(cfg *config.Configuration, db *sql.DB, logger *logger.Logger) *Container {
	c := &Container{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}

	c.initializeHandlers()
	return c
}

func (c *Container) initializeHandlers() {
	repositories := c.initRepositories()
	services := c.initServices(repositories)

	c.postHandler = handler.NewPostHandler(services.postService, c.logger)
	c.userHandler = handler.NewUserHandler(services.userService, c.logger)
	c.commentHandler = handler.NewCommentHandler(services.commentService, c.logger)
}

type repositories struct {
	userRepository    repository.UserRepository
	postRepository    repository.PostRepository
	commentRepository repository.CommentRepository
}

func (c *Container) initRepositories() repositories {
	return repositories{
		userRepository:    postgres.NewUserRepository(c.db),
		postRepository:    postgres.NewPostRepository(c.db),
		commentRepository: postgres.NewCommentRepository(c.db),
	}
}

type services struct {
	userService    service.UserService
	postService    service.PostService
	commentService service.CommentService
}

func (c *Container) initServices(r repositories) services {
	return services{
		userService:    application.NewUserService(r.userRepository),
		postService:    application.NewPostService(r.postRepository, r.commentRepository),
		commentService: application.NewCommentService(r.commentRepository),
	}
}
