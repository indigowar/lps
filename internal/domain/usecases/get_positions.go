package usecases

import (
	"context"
	"lps/internal/domain"
)

type GetPositionsUseCase interface {
	GetPositions(ctx context.Context) ([]domain.Position, error)
}
