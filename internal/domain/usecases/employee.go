package usecases

import (
	"context"
	"lps/internal/domain"

	"github.com/google/uuid"
)

type EmployeeUseCase interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Update(ctx context.Context, employee domain.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
}
