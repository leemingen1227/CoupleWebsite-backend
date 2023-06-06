package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/mail"
)

const (
	QueueCritical = "critical"
	QueueDefault = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendInvitationEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	return &RedisTaskProcessor{
		server: asynq.NewServer(
			redisOpt, 
			asynq.Config{
				Queues: map[string]int{
					QueueCritical: 10,
					QueueDefault: 5,
				},
			}),
		store: store,
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(TaskSendInvitationEmail, processor.ProcessTaskSendInvitationEmail)

	return processor.server.Start(mux)
}