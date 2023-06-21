package api

import (
	"database/sql"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/util"
	"github.com/leemingen1227/couple-server/worker"
	"github.com/lib/pq"
)

// createUserRequest is the request body for the CreateUser handler.
type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Name     string `json:"name" binding:"required,min=2,max=32"`
}

// userResponse is the response body for the CreateUser handler.
type userResponse struct {
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"password_changed_at"`
	CreateTime time.Time `json:"created_at"`
}

// new userResponse constructs a userResponse from a User.
func newUserResponse(user db.User) userResponse {
	return userResponse{
		Email:      user.Email,
		Name:       user.Name,
		UpdateTime: user.UpdateTime,
		CreateTime: user.CreateTime,
	}
}

// the handler for creating users.
// @Summary      SignUp
// @Description  Create a new user account
// @Tags         users
// @Param		 user	body	api.createUserRequest	true	"Create User"
// @Success      200  {object}  api.userResponse
// @Failure     400 
// @Router       /users/signup [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			ID:             uid,
			Email:          req.Email,
			PasswordDigest: hashedPassword,
			Name:           req.Name,
		},
		AfterCreate: func(user db.User) error {
			//Send verification email
			taskPayload := &worker.PayloadSendVerifyEmail{
				UserID: user.ID,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDirstributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(txResult.User)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

// the handler for login users.
// @Summary      Login
// @Description  Login to an user account
// @Tags         users
// @Param		 login_info	body	api.loginUserRequest	true	"Login User"
// @Success      200  {object}  api.loginUserResponse
// @Router       /users/login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.PasswordDigest)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//set the cookie
	// Set the session ID as a cookie in the response
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(session.ExpiresAt.Sub(time.Now()).Seconds()),
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // Set to true if using HTTPS
	}

	http.SetCookie(ctx.Writer, cookie)


	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)

}
