package auth

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	Login(ctx context.Context, login string, password string) (uuid.UUID, error)
}

func NewPostgresService(db *sqlx.DB) Service {
	return &postgresService{
		db: db,
	}
}

type postgresService struct {
	db *sqlx.DB
}

func (svc *postgresService) Login(ctx context.Context, login string, password string) (uuid.UUID, error) {
	var id uuid.UUID
	query := "SELECT employee FROM accounts WHERE login = $1 AND password = $2"

	var uuidArray []uuid.UUID
	if err := svc.db.Select(&uuidArray, query, login, password); err != nil {
		log.Println(err)
		return uuid.UUID{}, errors.New("error")
	}

	if len(uuidArray) > 0 {
		id = uuidArray[0]
	}

	return id, nil
}
