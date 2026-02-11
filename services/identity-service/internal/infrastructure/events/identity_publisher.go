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
	return p.publish(
		ctx,
		rabbitmq.IdentityUserRegisteredEvent,
		userID,
		rabbitmq.IdentityUserRegisteredEventPayload{
			UserID: userID,
			Email:  email,
		},
	)
}

func (p *IdentityEventPublisher) PublishEmailVerified(
	ctx context.Context,
	userID string,
	email string,
) error {
	return p.publish(
		ctx,
		rabbitmq.IdentityUserEmailVerifiedEvent,
		userID,
		rabbitmq.IdentityUserEmailVerifiedEventPayload{
			UserID: userID,
			Email:  email,
		},
	)
}

func (p *IdentityEventPublisher) PublishLogin(
	ctx context.Context,
	userID string,
	email string,
) error {
	return p.publish(
		ctx,
		rabbitmq.IdentityUserLoggedInEvent,
		userID,
		rabbitmq.IdentityUserLoggedInEventPayload{
			UserID: userID,
			Email:  email,
		},
	)
}

func (p *IdentityEventPublisher) publish(
	ctx context.Context,
	routingKey string,
	ownerID string,
	payload any,
) error {

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.rabbitmq.PublishMessage(
		ctx,
		routingKey,
		rabbitmq.AmqpMessage{
			OwnerID: ownerID,
			Data:    data,
		},
	)
}

