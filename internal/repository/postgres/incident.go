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

type incidentPostgres struct {
	db *sqlx.DB
}

// Delete implements IncidentUseCase.
func (u incidentPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := u.db.Exec("DELETE FROM incidents WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete incident")
	}
	return nil
}

// Get implements IncidentUseCase.
func (u incidentPostgres) Get(ctx context.Context, id uuid.UUID) (domain.Incident, error) {
	var i domain.Incident

	query := `
		SELECT id, employee, description, happenning_date
		FROM incidents WHERE id = $1
	`

	if err := u.db.Get(&i, query, id); err != nil {
		log.Println(err)
		return domain.Incident{}, fmt.Errorf("failed to load incident due to: %w", err)
	}

	return i, nil
}

// GetAll implements IncidentUseCase.
func (u incidentPostgres) GetAll(ctx context.Context) ([]domain.Incident, error) {
	query := `
		SELECT id, employee, description, happenning_date
		FROM incidents
	`
	incidents := make([]domain.Incident, 0)
	if err := u.db.Select(&incidents, query); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get incidents from postgres, due: %w", err)
	}
	return incidents, nil
}

// Update implements IncidentUseCase.
func (u incidentPostgres) Update(ctx context.Context, i domain.Incident) error {
	query := `
		UPDATE incidents
		SET employee = $1, description = $2, happenning_date = $3
		WHERE id = $4
	`

	if _, err := u.db.Exec(query, i.Employee, i.Description, i.Date, i.ID); err != nil {
		log.Println(err)
		return errors.New("failed to update incident")
	}

	return nil
}

func NewIncidentUseCase(db *sqlx.DB) usecases.IncidentUseCase {
	return incidentPostgres{
		db: db,
	}
}
