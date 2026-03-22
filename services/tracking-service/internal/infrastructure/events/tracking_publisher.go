package events

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
)

type TrackingEventPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewTrackingEventPublisher(rmq *rabbitmq.RabbitMQ) *TrackingEventPublisher {
	return &TrackingEventPublisher{
		rabbitmq: rmq,
	}
}
