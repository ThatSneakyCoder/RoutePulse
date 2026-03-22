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

			var err error

			switch msg.RoutingKey {

			case rabbitmq.IdentityUserRegisteredEvent:

				var payload rabbitmq.IdentityUserRegisteredEventPayload
				if e := json.Unmarshal(envelope.Data, &payload); e != nil {
					return e
				}

				err = c.service.InsertDomainEvent(
					ctx,
					"identity-service",
					rabbitmq.IdentityUserRegisteredEvent,
					payload.UserID,
					"",
				)

			case rabbitmq.IdentityUserEmailVerifiedEvent:

				var payload rabbitmq.IdentityUserEmailVerifiedEventPayload
				if e := json.Unmarshal(envelope.Data, &payload); e != nil {
					return e
				}

				err = c.service.InsertDomainEvent(
					ctx,
					"identity-service",
					rabbitmq.IdentityUserEmailVerifiedEvent,
					payload.UserID,
					"",
				)

			case rabbitmq.IdentityUserLoggedInEvent:

				var payload rabbitmq.IdentityUserLoggedInEventPayload
				if e := json.Unmarshal(envelope.Data, &payload); e != nil {
					return e
				}

				err = c.service.InsertDomainEvent(
					ctx,
					"identity-service",
					rabbitmq.IdentityUserLoggedInEvent,
					payload.UserID,
					"",
				)

			case rabbitmq.OrganizationCreatedEvent:

				var payload rabbitmq.OrganizationCreatedEventPayload
				if e := json.Unmarshal(envelope.Data, &payload); e != nil {
					return e
				}

				err = c.service.InsertDomainEvent(
					ctx,
					"organization-service",
					rabbitmq.OrganizationCreatedEvent,
					payload.OwnerUserID,
					payload.OrganizationID,
				)

			case rabbitmq.OrganizationMemberAddedEvent:

				var payload rabbitmq.OrganizationMemberAddedEventPayload
				if e := json.Unmarshal(envelope.Data, &payload); e != nil {
					return e
				}

				err = c.service.InsertDomainEvent(
					ctx,
					"organization-service",
					rabbitmq.OrganizationMemberAddedEvent,
					payload.UserID,
					payload.OrganizationID,
				)

			default:
				c.log.Infow("unknown routing key received",
					"routing_key", msg.RoutingKey,
				)
			}

			return err
		},
	)
}
