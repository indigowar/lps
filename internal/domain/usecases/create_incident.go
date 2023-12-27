package usecases

import (
	"context"

	"github.com/google/uuid"
)

type CreateIncidentUseCase interface {
	CreateIncident(ctx context.Context, user uuid.UUID, description string, date string) error
}
