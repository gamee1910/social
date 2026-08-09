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
	postHandler     *handler.PostHandler
	userHandler     *handler.UserHandler
	commentHandler  *handler.CommentHandler
	followerHandler *handler.FollowerHandler
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

func (c *Container) FollowerHandler() *handler.FollowerHandler {
	return c.followerHandler
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
	c.followerHandler = handler.NewFollowerHandler(services.followerService)
}

type repositories struct {
	userRepository     repository.UserRepository
	postRepository     repository.PostRepository
	commentRepository  repository.CommentRepository
	followerRepository repository.FollowerRepository
}

func (c *Container) initRepositories() repositories {
	return repositories{
		userRepository:     postgres.NewUserRepository(c.db, c.logger),
		postRepository:     postgres.NewPostRepository(c.db, c.logger),
		commentRepository:  postgres.NewCommentRepository(c.db, c.logger),
		followerRepository: postgres.NewFollowerRepository(c.db, c.logger),
	}
}

type services struct {
	userService     service.UserService
	postService     service.PostService
	commentService  service.CommentService
	followerService service.FollowerService
}

func (c *Container) initServices(r repositories) services {
	return services{
		userService:     application.NewUserService(r.userRepository, c.logger),
		postService:     application.NewPostService(r.postRepository, r.commentRepository, c.logger),
		commentService:  application.NewCommentService(r.commentRepository, c.logger),
		followerService: application.NewFollowerService(r.followerRepository, c.logger),
	}
}
