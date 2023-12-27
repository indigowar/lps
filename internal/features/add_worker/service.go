package addworker

import (
	"context"
	"fmt"
	"log"
	"lps/internal/domain"
	"lps/internal/domain/usecases"
	"lps/internal/repository/postgres"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	ConfirmUserCanAddWorker(ctx context.Context, id uuid.UUID) (bool, error)

	GetDepartments(ctx context.Context) ([]domain.Department, error)
	GetPositions(ctx context.Context) ([]domain.Position, error)

	CreateWorker(ctx context.Context, login string, surname string, name string, patronymic *string, phone string, position uuid.UUID, department uuid.UUID) error
}

type postgresService struct {
	db *sqlx.DB
}

// CreateWorker implements Service.
func (svc *postgresService) CreateWorker(ctx context.Context, login string, surname string, name string, patronymic *string, phone string, position uuid.UUID, department uuid.UUID) error {
	query := `SELECT create_worker($1, $2, $3, $4, $5, $6, $7)`
	_, err := svc.db.Exec(query, login, surname, name, patronymic, phone, position, department)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to create a worker in postgres, due to: %w", err)
	}
	return nil
}

// GetDepartments implements Service.
func (svc *postgresService) GetDepartments(ctx context.Context) ([]domain.Department, error) {
	return postgres.NewDepartmentUsecase(svc.db).GetAll(ctx)
}

// GetPositions implements Service.
func (svc *postgresService) GetPositions(ctx context.Context) ([]domain.Position, error) {
	return postgres.NewGetPositionsUseCase(svc.db).GetPositions(ctx)
}

// ConfirmUserCanAddWorker implements Service.
func (svc *postgresService) ConfirmUserCanAddWorker(ctx context.Context, id uuid.UUID) (bool, error) {
	role, err := postgres.NewGetUserRoleUseCase(svc.db).GetUserRole(ctx, id)
	if err != nil {
		return false, err
	}
	return role == usecases.UserRoleAdmin, nil
}

func NewPostgresService(db *sqlx.DB) Service {
	return &postgresService{
		db: db,
	}
}
