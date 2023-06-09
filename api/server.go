package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/token"
	"github.com/leemingen1227/couple-server/util"
	"github.com/leemingen1227/couple-server/worker"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "github.com/leemingen1227/couple-server/docs"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	config           util.Config
	store            db.Store
	tokenMaker       token.Maker
	router           *gin.Engine
	taskDirstributor worker.TaskDistributor
	redisClient 	*redis.Client
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor, redisClient *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:           config,
		store:            store,
		tokenMaker:       tokenMaker,
		taskDirstributor: taskDistributor,
		redisClient:	  redisClient,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.setupRouter()
	return server, nil
}

// @contact.name   Johnson Lee
// @title Couple Website API
// @version 1.0
// @description 
// @BasePath  /v1
// @host localhost:8080
func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Router := router.Group("/v1/")
	{
		userRouter := v1Router.Group("/users/")
		{
			userRouter.POST("/signup", server.createUser)
			userRouter.POST("/invitee_signup", server.invitedUserSignUp)	
			userRouter.POST("/login", server.loginUser)
		}

		verifyRouter := v1Router.Group("/verify/")
		{
			verifyRouter.GET("/verify_email", server.verifyEmail)
		}

		inviteRouter := v1Router.Group("/invite/")
		inviteRouter.Use(authMiddleware(server.tokenMaker, server.store, server.config, server.redisClient))
		{
			inviteRouter.POST("/", server.createInvitation)
		}

		blogRouter := v1Router.Group("/blogs/")
		blogRouter.Use(authMiddleware(server.tokenMaker, server.store, server.config, server.redisClient))
		{
			blogRouter.POST("/", server.createBlog)
			blogRouter.GET("/blog/:blogID", server.getBlogByBlogID)
			blogRouter.GET("/:pairID", server.getBlogsByPairID)
		}
	}


	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
