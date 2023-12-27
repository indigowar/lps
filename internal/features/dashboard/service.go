package dashboard

import (
	"context"
	"lps/internal/domain"
	"lps/internal/repository/postgres"

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

	GetPositions(ctx context.Context) ([]domain.Position, error)
}

type postgresService struct {
	db *sqlx.DB
}

// GetPositions implements Service.
func (svc *postgresService) GetPositions(ctx context.Context) ([]domain.Position, error) {
	return postgres.NewGetPositionsUseCase(svc.db).GetPositions(ctx)
}

func (*postgresService) GetUserRole(ctx context.Context, id uuid.UUID) (UserRole, error) {
	return UserRoleAdmin, nil
}

func NewPostgrseService(db *sqlx.DB) Service {
	return &postgresService{
		db: db,
	}
}
