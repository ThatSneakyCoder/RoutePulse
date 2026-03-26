package events

import (
	"context"
	"encoding/json"

	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
)

type APIGatewayEventPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewAPIGatewayEventPublisher(rmq *rabbitmq.RabbitMQ) *APIGatewayEventPublisher {
	return &APIGatewayEventPublisher{
		rabbitmq: rmq,
	}
}

func (p *APIGatewayEventPublisher) PublishTrackingLocationUpdated(
	ctx context.Context,
	payload rabbitmq.TrackingDriverLocationUpdatedEventPayload,
) error {
	return p.publish(
		ctx,
		rabbitmq.TrackingDriverLocationUpdatedEvent,
		payload.DriverID,
		payload,
	)
}

func (p *APIGatewayEventPublisher) publish(
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
