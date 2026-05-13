package services

import (
	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/rs/zerolog"
)

func logAccessDenied(logger zerolog.Logger, actor dto.Actor, resource string, resourceID string, message string) {
	event := logger.Warn().
		Str("actor_id", actor.ID).
		Str("actor_role", string(actor.Role)).
		Str("resource", resource)
	if resourceID != "" {
		event = event.Str("resource_id", resourceID)
	}

	event.Msg(message)
}

func logRepositoryError(logger zerolog.Logger, err error, resource string, resourceID string, message string) {
	event := logger.Error().
		Err(err).
		Str("resource", resource)
	if resourceID != "" {
		event = event.Str("resource_id", resourceID)
	}

	event.Msg(message)
}
