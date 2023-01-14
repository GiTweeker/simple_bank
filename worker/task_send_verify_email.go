package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (r RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context,
	payload *PayloadSendVerifyEmail, opt ...asynq.Option) error {

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opt...)

	info, err := r.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enquire Task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail

	json.Unmarshal(task.Payload(), &payload)

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarhsall payload : %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exist %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")

	return nil
}
