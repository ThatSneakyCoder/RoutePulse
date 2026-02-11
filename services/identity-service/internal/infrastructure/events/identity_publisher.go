package events

import (
	"context"
	"encoding/json"

	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
)

type IdentityEventPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewIdentityEventPublisher(rmq *rabbitmq.RabbitMQ) *IdentityEventPublisher {
	return &IdentityEventPublisher{
		rabbitmq: rmq,
	}
}

// PublishUserRegistered emits an event when a new user is registered
func (p *IdentityEventPublisher) PublishUserRegistered(
	ctx context.Context,
	userID string,
	email string,
) error {

	payload := map[string]any{
		"user_id": userID,
		"email":   email,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.rabbitmq.PublishMessage(
		ctx,
		rabbitmq.IdentityUserRegisteredEvent,
		rabbitmq.AmqpMessage{
			OwnerID: userID,
			Data:    data,
		},
	)
}
