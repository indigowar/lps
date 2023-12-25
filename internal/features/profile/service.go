package profile

import (
	"context"
	"errors"
	"log"
	"lps/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	GetUserInfo(ctx context.Context, id uuid.UUID) (domain.Account, domain.Employee, error)

	UpdateAccountLogin(ctx context.Context, id uuid.UUID, password string, login string) error
	UpdateAccountPassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) error

	UpdateAccount(ctx context.Context, id uuid.UUID, password string, newLogin string, newPassword string) error
}

type postgresService struct {
	db *sqlx.DB
}

// GetUserInfo implements Service.
func (svc *postgresService) GetUserInfo(ctx context.Context, id uuid.UUID) (domain.Account, domain.Employee, error) {
	var e domain.Employee
	var a domain.Account

	query := `
		SELECT
			s.id, s.surname, s.name, s.patronymic,
			s.phone_number, s.position, s.department,
			a.login, a.password, a.activated, a.employee
		FROM staff s
		LEFT JOIN accounts a ON s.id = a.employee
		WHERE s.id = $1
	`

	err := svc.db.QueryRow(query, id).Scan(&e.ID, &e.Surname, &e.Name, &e.Patronymic, &e.PhoneNumber, &e.Position, &e.Department, &a.Login, &a.Password, &a.IsActivated, &a.Employee)
	if err != nil {
		log.Println(err)
		return domain.Account{}, domain.Employee{}, errors.New("failed to load user info")
	}
	return a, e, nil
}

// UpdateAccountLogin implements Service.
func (svc *postgresService) UpdateAccountLogin(ctx context.Context, id uuid.UUID, password string, login string) error {
	if _, err := svc.db.Exec("UPDATE accounts SET login = $1 WHERE id = $3 AND password = $2", login, password, id); err != nil {
		log.Print(err)
		return errors.New("failed to update account")
	}
	return nil
}

// UpdateAccountPassword implements Service.
func (svc *postgresService) UpdateAccountPassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) error {
	if _, err := svc.db.Exec("UPDATE accounts SET password = $1 WHERE id = $3 AND password = $2", newPassword, oldPassword, id); err != nil {
		log.Print(err)
		return errors.New("failed to update account")
	}
	return nil
}

func (svc *postgresService) UpdateAccount(ctx context.Context, id uuid.UUID, password string, newLogin string, newPassword string) error {
	tx := svc.db.MustBegin()
	if newLogin != "" {
		tx.MustExec("UPDATE accounts SET login = $1 WHERE employee = $2 AND password = $3", newLogin, id, password)
	}
	if newPassword != "" {
		tx.MustExec("UPDATE accounts SET password = $1 WHERE employee = $2 AND password = $3", newPassword, id, password)
	}

	if err := tx.Commit(); err != nil {
		log.Print(err)
		return errors.New("failed to update")
	}
	return nil
}

func NewPostgresService(db *sqlx.DB) Service {
	return &postgresService{
		db: db,
	}
}
