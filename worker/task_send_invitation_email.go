package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendInvitationEmail = "task:send_invitation_email"

type PayloadSendInvitationEmail struct {
	InvitationID int64 `json:"invitation_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendInvitationEmail(
	ctx context.Context,
	payload *PayloadSendInvitationEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendInvitationEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).
	Int("max_retry", info.MaxRetry).Msg("task enqueued")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendInvitationEmail (ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendInvitationEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	invitation, err := processor.store.GetInvitation(ctx, payload.InvitationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invitation doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get invitation: %w", err)
	}

	user, err := processor.store.GetUser(ctx, invitation.InviterID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	subject := "You have been invited to join Couple Website!"
	verifyUrl := fmt.Sprintf("http://localhost:8080/invitee_signup?invitation_id=%d&invitation_token=%s", invitation.ID, invitation.InvitationToken)
	content := fmt.Sprintf(`Hello,<br/>
	You are invited to join Couple Website by %s.<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.Name, verifyUrl)
	to := []string{invitation.InviteeEmail}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
	Str("invitee_email", invitation.InviteeEmail).Msg("task processed")

	return nil

}