package service

import (
	"context"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
)

func (s *AnalyticsService) InsertIdentityUserRegistered(
	ctx context.Context,
	event rabbitmq.IdentityUserRegisteredEventPayload,
) error {

	s.log.Infow("handling identity.user.registered event",
		"user_id", event.UserID,
		"email", event.Email,
	)

	err := s.repo.InsertIdentityUserRegistered(
		ctx,
		time.Now(),
		event.UserID,
		event.Email,
	)

	if err != nil {
		s.log.Errorw("failed to insert identity registration event",
			"user_id", event.UserID,
			"err", err,
		)
		return err
	}

	s.log.Infow("identity registration event inserted successfully",
		"user_id", event.UserID,
	)

	return nil
}
