package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	Login(ctx context.Context, login string, password string) (uuid.UUID, error)

	HasUnactivatedUser(ctx context.Context, login string) (bool, error)

	ActivateUser(ctx context.Context, login string, password string) (uuid.UUID, error)
}

func NewPostgresService(db *sqlx.DB) Service {
	return &postgresService{
		db: db,
	}
}

type postgresService struct {
	db *sqlx.DB
}

// ActivateUser implements Service.
func (svc *postgresService) ActivateUser(ctx context.Context, login string, password string) (uuid.UUID, error) {
	id := make([]uuid.UUID, 0)
	if err := svc.db.Select(&id, `SELECT activate_account($1, $2)`, login, password); err != nil {
		log.Println(err)
		return uuid.UUID{}, fmt.Errorf("failed to activate account, due to: %w", err)
	}
	return id[0], nil
}

// HasUnactivatedUser implements Service.
func (svc *postgresService) HasUnactivatedUser(ctx context.Context, login string) (bool, error) {
	query := `SELECT activated FROM accounts WHERE login = $1`
	result := make([]bool, 0)
	if err := svc.db.Select(&result, query, login); err != nil {
		log.Println(err)
		return false, errors.New("failed to load value")
	}
	if len(result) == 0 {
		return false, errors.New("user does not exist")
	}
	return !result[0], nil
}

func (svc *postgresService) Login(ctx context.Context, login string, password string) (uuid.UUID, error) {
	var id uuid.UUID
	query := "SELECT employee FROM accounts WHERE login = $1 AND password = $2"

	var uuidArray []uuid.UUID
	if err := svc.db.Select(&uuidArray, query, login, password); err != nil {
		log.Println(err)
		return uuid.UUID{}, errors.New("error")
	}

	if len(uuidArray) == 0 {
		return uuid.UUID{}, errors.New("user not found")
	}

	if len(uuidArray) > 0 {
		id = uuidArray[0]
	}

	return id, nil
}
