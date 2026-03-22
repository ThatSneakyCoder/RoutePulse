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

// --------------------
// Events
// --------------------

// User Registered (no org yet)
func (p *IdentityEventPublisher) PublishUserRegistered(
	ctx context.Context,
	userID string,
) error {
	return p.publish(
		ctx,
		rabbitmq.IdentityUserRegisteredEvent,
		userID,
		rabbitmq.IdentityUserRegisteredEventPayload{
			UserID: userID,
		},
	)
}

// Email Verified (no org yet)
func (p *IdentityEventPublisher) PublishEmailVerified(
	ctx context.Context,
	userID string,
) error {
	return p.publish(
		ctx,
		rabbitmq.IdentityUserEmailVerifiedEvent,
		userID,
		rabbitmq.IdentityUserEmailVerifiedEventPayload{
			UserID: userID,
		},
	)
}

// Login (org-aware)
func (p *IdentityEventPublisher) PublishLogin(
	ctx context.Context,
	userID string,
	orgID string,
) error {
	return p.publish(
		ctx,
		rabbitmq.IdentityUserLoggedInEvent,
		userID,
		rabbitmq.IdentityUserLoggedInEventPayload{
			UserID: userID,
			OrgID:  orgID,
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
