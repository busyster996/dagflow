package service

import (
	"context"

	"github.com/busyster996/dagflow/internal/pubsub"
)

type SEventService struct {
}

func Event() *SEventService {
	return &SEventService{}
}

func (e *SEventService) Subscribe(ctx context.Context, event chan string) error {
	return pubsub.SubscribeEvent(ctx, func(data string) {
		event <- data
	})
}
