package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"lps/internal/domain"
	"lps/internal/domain/usecases"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresPosition struct {
	db *sqlx.DB
}

// Delete implements usecases.PositionUseCase.
func (u *postgresPosition) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM positions WHERE id = $1`
	if _, err := u.db.Exec(query, id); err != nil {
		log.Println(err)
		return errors.New("failed to delete")
	}
	return nil
}

// Get implements usecases.PositionUseCase.
func (u *postgresPosition) Get(ctx context.Context, id uuid.UUID) (domain.Position, error) {
	d := make([]domain.Position, 0)
	if err := u.db.Select(&d, "SELECT id, title, level FROM positions WHERE id = $1", id); err != nil {
		log.Println(err)
		return domain.Position{}, fmt.Errorf("failed to load department due to: %w", err)
	}

	if len(d) < 1 {
		return domain.Position{}, errors.New("not found")
	}

	return d[0], nil
}

// GetAll implements usecases.PositionUseCase.
func (u *postgresPosition) GetAll(ctx context.Context) ([]domain.Position, error) {
	query := `
	SELECT id, title, level 
	FROM positions
	`
	positions := make([]domain.Position, 0)
	if err := u.db.Select(&positions, query); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get positions from postgres, due: %w", err)
	}
	return positions, nil
}

// Update implements usecases.PositionUseCase.
func (u *postgresPosition) Update(ctx context.Context, id uuid.UUID, name string, level string) error {
	_, err := u.db.Exec(`UPDATE positions SET title = $1, level = $2 WHERE id = $3`, name, level, id)
	if err != nil {
		log.Println(err)
		return errors.New("failed to update")
	}
	return nil
}

func NewPositionUseCase(db *sqlx.DB) usecases.PositionUseCase {
	return &postgresPosition{
		db: db,
	}
}
