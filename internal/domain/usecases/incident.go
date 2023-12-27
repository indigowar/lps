package usecases

import (
	"context"
	"lps/internal/domain"

	"github.com/google/uuid"
)

type IncidentUseCase interface {
	GetAll(ctx context.Context) ([]domain.Incident, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Incident, error)
	Update(ctx context.Context, incident domain.Incident) error
	Delete(ctx context.Context, id uuid.UUID) error
}
