package events

import (
	"context"
	"encoding/json"

	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
)

type OrganizationEventPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewOrganizationEventPublisher(rmq *rabbitmq.RabbitMQ) *OrganizationEventPublisher {
	return &OrganizationEventPublisher{
		rabbitmq: rmq,
	}
}

// --------------------
// Events
// --------------------

// Organization Created
func (p *OrganizationEventPublisher) PublishOrganizationCreated(
	ctx context.Context,
	orgID string,
	ownerUserID string,
) error {
	return p.publish(
		ctx,
		rabbitmq.OrganizationCreatedEvent,
		ownerUserID,
		rabbitmq.OrganizationCreatedEventPayload{
			OrganizationID: orgID,
			OwnerUserID:    ownerUserID,
		},
	)
}

// Member Added
func (p *OrganizationEventPublisher) PublishOrganizationMemberAdded(
	ctx context.Context,
	userID string,
	orgID string,
	role string,
) error {
	return p.publish(
		ctx,
		rabbitmq.OrganizationMemberAddedEvent,
		userID,
		rabbitmq.OrganizationMemberAddedEventPayload{
			UserID:         userID,
			OrganizationID: orgID,
			Role:           role,
		},
	)
}

func (p *OrganizationEventPublisher) publish(
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
