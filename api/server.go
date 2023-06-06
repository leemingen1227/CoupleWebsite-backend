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
)

type Server struct {
	config           util.Config
	store            db.Store
	tokenMaker       token.Maker
	router           *gin.Engine
	taskDirstributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:           config,
		store:            store,
		tokenMaker:       tokenMaker,
		taskDirstributor: taskDistributor,
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
		inviteRouter.Use(authMiddleware(server.tokenMaker))
		{
			inviteRouter.POST("/", server.createInvitation)
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
