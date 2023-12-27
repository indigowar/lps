package postgres

import (
	"context"
	"errors"
	"log"
	"lps/internal/domain/usecases"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresCreateIncident struct {
	db *sqlx.DB
}

// CreateIncident implements usecases.CreateIncidentUseCase.
func (u *postgresCreateIncident) CreateIncident(ctx context.Context, user uuid.UUID, description string, date string) error {
	query := `INSERT INTO incidents(employee, description, happenning_date) VALUES ($1, $2, $3)`

	if _, err := u.db.Exec(query, user, description, date); err != nil {
		log.Println(err)
		return errors.New("failed to insert incidents")
	}
	return nil
}

func NewCreateIncidentsUseCase(db *sqlx.DB) usecases.CreateIncidentUseCase {
	return &postgresCreateIncident{
		db: db,
	}
}
