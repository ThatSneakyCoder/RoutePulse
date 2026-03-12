package events

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
)

type FleetEventPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewFleetEventPublisher(rmq *rabbitmq.RabbitMQ) *FleetEventPublisher {
	return &FleetEventPublisher{
		rabbitmq: rmq,
	}
}
