package usecase

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/usecase/port/queue"
	"context"
)

type MediaEncodeStatusUseCase struct {
	queue queue.QueueHandler
}

func NewMediaEncodeStatusUseCase(q queue.QueueHandler) *MediaEncodeStatusUseCase {
	return &MediaEncodeStatusUseCase{
		queue: q,
	}
}

type MediaEncodeStatusInput struct {
	ReqCtx         context.Context
	OwnerAccountId *value_object.OwnerAccountId
	WorkerFunc     func(map[string]interface{})
}

func (u *MediaEncodeStatusUseCase) Execute(input *MediaEncodeStatusInput) error {
	return u.queue.Status(input.ReqCtx, input.OwnerAccountId, input.WorkerFunc)
}
