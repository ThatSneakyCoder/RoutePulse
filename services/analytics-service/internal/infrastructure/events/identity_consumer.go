package events

import (
	"context"
	"encoding/json"

	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/service"
	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type identityConsumer struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  *service.AnalyticsService
	log      *zap.SugaredLogger
}

func NewIdentityConsumer(rabbitmq *rabbitmq.RabbitMQ, service *service.AnalyticsService, log *zap.SugaredLogger) *identityConsumer {
	return &identityConsumer{
		rabbitmq: rabbitmq,
		service:  service,
		log:      log,
	}
}

func (c *identityConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessages(
		rabbitmq.AnalyticsIdentityQueue,
		func(ctx context.Context, msg amqp091.Delivery) error {

			c.log.Infow("received message from analytics.identity queue",
				"routing_key", msg.RoutingKey,
			)

			var envelope rabbitmq.AmqpMessage
			if err := json.Unmarshal(msg.Body, &envelope); err != nil {
				c.log.Errorw("failed to unmarshal envelope", "err", err)
				return err
			}

			if msg.RoutingKey != rabbitmq.IdentityUserRegisteredEvent {
				c.log.Infow("ignoring unrelated routing key",
					"routing_key", msg.RoutingKey,
				)
				return nil
			}

			var payload rabbitmq.IdentityUserRegisteredEventPayload
			if err := json.Unmarshal(envelope.Data, &payload); err != nil {
				c.log.Errorw("failed to unmarshal payload", "err", err)
				return err
			}

			return c.service.InsertIdentityUserRegistered(ctx, rabbitmq.IdentityUserRegisteredEventPayload{
				UserID: payload.UserID,
				Email:  payload.Email,
			})
		},
	)
}
