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

type staffPostgres struct {
	db *sqlx.DB
}

// Delete implements usecases.EmployeeUseCase.
func (u staffPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := u.db.Exec("DELETE FROM staff WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete")
	}
	return nil
}

// Get implements usecases.EmployeeUseCase.
func (u staffPostgres) Get(ctx context.Context, id uuid.UUID) (domain.Employee, error) {
	d := make([]domain.Employee, 0)

	query := `
		SELECT id, surname, name, patronymic, phone_number, position, department
		FROM staff WHERE id = $1
	`

	if err := u.db.Select(&d, query, id); err != nil {
		log.Println(err)
		return domain.Employee{}, fmt.Errorf("failed to load employee due to: %w", err)
	}

	if len(d) < 1 {
		return domain.Employee{}, errors.New("not found")
	}

	return d[0], nil
}

// GetAll implements usecases.EmployeeUseCase.
func (u staffPostgres) GetAll(ctx context.Context) ([]domain.Employee, error) {
	query := `
		SELECT id, surname, name, patronymic, phone_number, position, department
		FROM staff	
	`
	staff := make([]domain.Employee, 0)
	if err := u.db.Select(&staff, query); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get departments from postgres, due: %w", err)
	}
	return staff, nil
}

// Update implements usecases.EmployeeUseCase.
func (u staffPostgres) Update(ctx context.Context, e domain.Employee) error {
	query := `
		UPDATE staff
		SET surname = $1, name = $2, patronymic = $3,
			phone_number = $4, position = $5, department = $6
		WHERE id = $7
	`

	if _, err := u.db.Exec(query, e.Surname, e.Name, e.Patronymic, e.PhoneNumber, e.Position, e.Department, e.ID); err != nil {
		log.Println(err)
		return errors.New("failed to update")
	}

	return nil
}

func NewEmployeeUseCase(db *sqlx.DB) usecases.EmployeeUseCase {
	return staffPostgres{
		db: db,
	}
}
