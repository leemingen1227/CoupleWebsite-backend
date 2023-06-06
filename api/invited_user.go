package api

import (
	"time"

	"net/http"

	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/leemingen1227/couple-server/util"
	
)

// The process of invited user signing up:
// 1. Invited user click the link in the invitation email
// 2. Invited user fill in the sign up form, providing email, password, name
// 3. Invited user click the sign up button
// 4. Update the invitation record in the database, set is_accepted to true
// 5. Create a new Pair record in the database
// 6. Update the invitee user record in the database, set pair_id to the new pair id
// 7. Create the invited user record in the database, with the provided email, password, name, pair_id
// 8. Create two UserPair records in the database, one for the inviter, one for the invitee
// We don't have to verify the invited user's email address, because the invitation email is sent to the invited user's email address.

type invitedUserSignUpRequest struct {
	InvitationID    int64  `json:"invitation_id"`
	InvitationToken string `json:"invitation_token"`
	Password        string `json:"password"`
	Name            string `json:"name"`
}


type invitedUserSignUpResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	CreateTime time.Time `json:"create_at"`
}

func newInvitedUserSignupResponse (user db.User) invitedUserSignUpResponse {
	return invitedUserSignUpResponse{
		Email: user.Email,
		Name: user.Name,
		CreateTime: user.CreateTime,
	}
}

// Invited user sign up
//the handler for invitee to sign up
// @Summary      Invitee SignUp
// @Description  for invitee to sign up
// @Tags         users
// @Param		 signup_info	body	api.invitedUserSignUpRequest	true	"Create User"
// @Success      200  {object}  api.invitedUserSignUpResponse
// @Router       /users/invitee_signup [post]
func (server *Server) invitedUserSignUp(ctx *gin.Context) {
	var req invitedUserSignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	invitation, err := server.store.GetInvitation(ctx, req.InvitationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.InvitedUserSignUpTxParams{
		InvitationID: req.InvitationID,
		InvitationToken: req.InvitationToken,
		CreateUserParams: db.CreateUserParams{
			Email: invitation.InviteeEmail,
			PasswordDigest: hashedPassword,
			Name: req.Name,
			IsEmailVerified: true,
		},
	}

	var txResult db.InvitedUserSignUpTxResult

	txResult, err = server.store.InvitedUserSignUpTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newInvitedUserSignupResponse(txResult.User)
	ctx.JSON(http.StatusOK, rsp)
}
