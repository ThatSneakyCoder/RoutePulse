package events

import (
	"context"
	"encoding/json"

	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/service"
	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type trackingConsumer struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  *service.TrackingService
	log      *zap.SugaredLogger
}

func NewTrackingConsumer(rabbitmq *rabbitmq.RabbitMQ, service *service.TrackingService, log *zap.SugaredLogger) *trackingConsumer {
	return &trackingConsumer{
		rabbitmq: rabbitmq,
		service:  service,
		log:      log,
	}
}

func (c *trackingConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessages(
		rabbitmq.TrackingLocationQueue,
		func(ctx context.Context, msg amqp091.Delivery) error {
			c.log.Infow("received message from tracking.location queue", "routing_key", msg.RoutingKey)

			var envelope rabbitmq.AmqpMessage
			if err := json.Unmarshal(msg.Body, &envelope); err != nil {
				c.log.Errorw("failed to unmarshal tracking message envelope", "err", err)
				return err
			}

			switch msg.RoutingKey {
			case rabbitmq.TrackingDriverLocationUpdatedEvent:
				var payload rabbitmq.TrackingDriverLocationUpdatedEventPayload
				if err := json.Unmarshal(envelope.Data, &payload); err != nil {
					c.log.Errorw("failed to unmarshal tracking location payload", "err", err)
					return err
				}

				return c.service.StoreLocationUpdate(ctx, domain.TrackingLocationUpdate{
					TripID:     payload.TripID,
					DriverID:   payload.DriverID,
					VehicleID:  payload.VehicleID,
					Latitude:   payload.Latitude,
					Longitude:  payload.Longitude,
					RecordedAt: payload.RecordedAt,
					Sequence:   payload.Sequence,
				})
			default:
				c.log.Infow("unknown tracking routing key received", "routing_key", msg.RoutingKey)
				return nil
			}
		},
	)
}
