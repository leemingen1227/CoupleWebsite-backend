package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/token"
	"github.com/leemingen1227/couple-server/util"
	"github.com/leemingen1227/couple-server/worker"
)

type InviteRequest struct {
	InviterID   uuid.UUID `json:"inviter_id"`
	InviteeEmail string `json:"invitee_email"`
}

type InviteResponse struct {
	ID int64 `json:"id"`
	InviterID uuid.UUID `json:"inviter_id"`
	InviteeEmail string `json:"invitee_email"`
	InvitationToken string `json:"invitation_token"`
	IsAccepted bool `json:"is_accepted"`
	CreateTime time.Time `json:"create_time"`
}

func newInviteResponse(invitation db.Invitation) InviteResponse {
	return InviteResponse{
		ID: invitation.ID,
		InviterID: invitation.InviterID,
		InviteeEmail: invitation.InviteeEmail,
		InvitationToken: invitation.InvitationToken,
		IsAccepted: invitation.IsAccepted,
		CreateTime: invitation.CreateTime,
	}
}

//the handler for creating inviation
// @Summary      Invite
// @Description  Invite new user to create a pair
// @Tags         invite
// @Param        Authorization     header    string     true   "Bearer token"
// @Param		 invite_info	body	api.InviteRequest	true	"Invite User"
// @Success      200  {object}  api.InviteResponse
// @Router       /invite [post]
func (server *Server) createInvitation(ctx *gin.Context) {
	var req InviteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Get the inviter id from the context 
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	req.InviterID = authPayload.UserID

	//check if inviter has verified email
	inviter, err := server.store.GetUser(ctx, req.InviterID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !inviter.IsEmailVerified {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateInvitationTxParams{
		CreateInvitationParams: db.CreateInvitationParams{
			InviterID: req.InviterID,
			InviteeEmail: req.InviteeEmail,
			InvitationToken: util.RandomString(32),
		},
		AfterCreate: func(invitation db.Invitation) error {
			//send email to invitee
			taskPayload := &worker.PayloadSendInvitationEmail{
				InvitationID: invitation.ID,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.Timeout(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDirstributor.DistributeTaskSendInvitationEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateInvitationTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	invitation := txResult.Invitation
	rsp := newInviteResponse(invitation)
	ctx.JSON(http.StatusOK, rsp)

}