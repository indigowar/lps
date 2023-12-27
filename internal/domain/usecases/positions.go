package usecases

import (
	"context"
	"lps/internal/domain"

	"github.com/google/uuid"
)

type PositionUseCase interface {
	GetAll(ctx context.Context) ([]domain.Position, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Position, error)
	Update(ctx context.Context, id uuid.UUID, name string, level string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
