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

type postgresDepartment struct {
	db *sqlx.DB
}

// Delete implements usecases.DepartmentUseCase.
func (u *postgresDepartment) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM departments WHERE id = $1`
	if _, err := u.db.Exec(query, id); err != nil {
		log.Println(err)
		return errors.New("failed to delete")
	}
	return nil
}

// Get implements usecases.DepartmentUseCase.
func (u *postgresDepartment) Get(ctx context.Context, id uuid.UUID) (domain.Department, error) {
	d := make([]domain.Department, 0)

	if err := u.db.Select(&d, "SELECT id, name FROM departments WHERE id = $1", id); err != nil {
		log.Println(err)
		return domain.Department{}, fmt.Errorf("failed to load department due to: %w", err)
	}

	if len(d) < 1 {
		return domain.Department{}, errors.New("not found")
	}

	return d[0], nil
}

// GetAll implements usecases.DepartmentUseCase.
func (u *postgresDepartment) GetAll(ctx context.Context) ([]domain.Department, error) {
	query := `
	SELECT id, name
	FROM departments
	`
	deparmentes := make([]domain.Department, 0)
	if err := u.db.Select(&deparmentes, query); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get departments from postgres, due: %w", err)
	}
	return deparmentes, nil
}

// Update implements usecases.DepartmentUseCase.
func (u *postgresDepartment) Update(ctx context.Context, id uuid.UUID, name string) error {
	_, err := u.db.Exec(`UPDATE departments SET name = $1 WHERE id = $2`, name, id)
	if err != nil {
		log.Println(err)
		return errors.New("failed to update")
	}

	return nil
}

func NewDepartmentUsecase(db *sqlx.DB) usecases.DepartmentUseCase {
	return &postgresDepartment{
		db: db,
	}
}
