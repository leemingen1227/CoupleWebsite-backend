package worker

import (
	"github.com/hibiken/asynq"
	"context"
)
type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
	DistributeTaskSendInvitationEmail(
		ctx context.Context,
		payload *PayloadSendInvitationEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	return &RedisTaskDistributor{
		client: asynq.NewClient(redisOpt),
	}
}