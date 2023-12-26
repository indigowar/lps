package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"lps/internal/domain"
	"lps/internal/domain/usecases"
)

type postgresGetDepartments struct {
	db *sqlx.DB
}

// GetDepartments implements usecases.GetDepartmentsUseCase.
func (usecase *postgresGetDepartments) GetDepartments(ctx context.Context) ([]domain.Department, error) {
	query := `
	SELECT id, name
	FROM departments
	`
	deparmentes := make([]domain.Department, 0)
	if err := usecase.db.Select(&deparmentes, query); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get departments from postgres, due: %w", err)
	}
	return deparmentes, nil
}

func NewGetDepartmentsUseCase(db *sqlx.DB) usecases.GetDepartmentsUseCase {
	return &postgresGetDepartments{
		db: db,
	}
}
