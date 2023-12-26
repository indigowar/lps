package usecases

import (
	"context"
	"lps/internal/domain"
)

type GetDepartmentsUseCase interface {
	GetDepartments(ctx context.Context) ([]domain.Department, error)
}
