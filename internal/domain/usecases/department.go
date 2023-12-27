package usecases

import (
	"context"
	"lps/internal/domain"

	"github.com/google/uuid"
)

type DepartmentUseCase interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Department, error)
	GetAll(ctx context.Context) ([]domain.Department, error)
	Update(ctx context.Context, id uuid.UUID, name string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
