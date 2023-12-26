package dashboard

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRole int

const (
	UserRoleAdmin UserRole = iota
	UserRoleHead
	UserRoleStaff
)

type Service interface {
	GetUserRole(ctx context.Context, id uuid.UUID) (UserRole, error)
}

type postgresService struct {
	db *sqlx.DB
}

func (*postgresService) GetUserRole(ctx context.Context, id uuid.UUID) (UserRole, error) {
	return UserRoleAdmin, nil
}

func NewPostgrseService(db *sqlx.DB) Service {
	return &postgresService{
		db: db,
	}
}
