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

func (p *OrganizationEventPublisher) PublishOrganizationCreated(
	ctx context.Context,
	orgName string,
	orgID string,
	userID string,
) error {
	return p.publish(
		ctx,
		rabbitmq.OrganizationCreatedEvent,
		userID,
		rabbitmq.OrganizationCreatedEventPayload{
			OrganizationID: orgID,
			Name: orgName,
			OwnerUserID: userID,
		},
	)
}

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
			UserID: userID,
			OrganizationID: orgID,
			Role: role,
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

