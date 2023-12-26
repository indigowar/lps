package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"lps/internal/domain"
	"lps/internal/domain/usecases"
)

type getPositionsPostgres struct {
	db *sqlx.DB
}

// GetPositions implements usecases.GetPositionsUseCase.
func (usecase *getPositionsPostgres) GetPositions(ctx context.Context) ([]domain.Position, error) {
	query := `
		SELECT id, title, level
		FROM positions
	`

	positions := make([]domain.Position, 0)

	if err := usecase.db.Select(&positions, query); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get positions from postgres, due to: %w", err)
	}

	return positions, nil
}

func NewGetPositionsUseCase(db *sqlx.DB) usecases.GetPositionsUseCase {
	return &getPositionsPostgres{
		db: db,
	}
}
