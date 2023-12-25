package profile

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserInfo struct {
	ID          uuid.UUID
	Surname     string
	Name        string
	Patronymic  *string
	PhoneNumber string
	Postion     string
	Department  string
	Login       string
	Password    string
	Activated   bool
}

type Service interface {
	GetUserInfo(ctx context.Context, id uuid.UUID) (UserInfo, error)

	UpdateAccountLogin(ctx context.Context, id uuid.UUID, password string, login string) error
	UpdateAccountPassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) error

	UpdateAccount(ctx context.Context, id uuid.UUID, password string, newLogin string, newPassword string) error
	UpdateEmployee(ctx context.Context, id uuid.UUID, surname string, name string, patronymic *string, phone_number string) error
}

type postgresService struct {
	db *sqlx.DB
}

// GetUserInfo implements Service.
func (svc *postgresService) GetUserInfo(ctx context.Context, id uuid.UUID) (UserInfo, error) {
	var d UserInfo

	query := `
		SELECT
			id, surname, name, patronymic,
			phone_number, position, department,
			login, password, activated
		FROM employee_details_view
		WHERE id = $1	
	`

	err := svc.db.QueryRow(query, id).Scan(&d.ID, &d.Surname, &d.Name, &d.Patronymic, &d.PhoneNumber, &d.Postion, &d.Department, &d.Login, &d.Password, &d.Activated)
	if err != nil {
		log.Println(err)
		return UserInfo{}, errors.New("failed to retrieve info about user")
	}
	return d, nil
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

func (svc *postgresService) UpdateEmployee(ctx context.Context, id uuid.UUID, surname string, name string, patronymic *string, phone_number string) error {
	tx := svc.db.MustBegin()

	if surname != "" {
		tx.Exec("UPDATE staff SET surname = $1 WHERE id = $2", surname, id)
	}

	if name != "" {
		tx.Exec("UPDATE staff SET name = $1 WHERE id = $2", name, id)
	}

	if patronymic != nil {
		if *patronymic != "" {
			tx.Exec("UPDATE staff SET patronymic = $1 WHERE id = $2", *patronymic, id)
		}
	}

	if phone_number != "" {
		tx.Exec("UPDATE staff SET phone_number = $1 WHERE id = $2", phone_number, id)
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return errors.New("failed to save")
	}
	return nil
}

func NewPostgresService(db *sqlx.DB) Service {
	return &postgresService{
		db: db,
	}
}
