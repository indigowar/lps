package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"lps/internal/domain/usecases"
)

type getUserRole struct {
	db *sqlx.DB
}

func (g *getUserRole) GetUserRole(ctx context.Context, id uuid.UUID) (usecases.UserRole, error) {
	query := `
		SELECT p.level
		FROM positions p
		RIGHT JOIN staff s ON s.position = p.id
		WHERE s.id = $1	
	`
	levels := make([]string, 0)
	if err := g.db.Select(&levels, query, id); err != nil {
		log.Println(err)
		return usecases.UserRoleNone, fmt.Errorf("failed to get user role from postgres, due: %w", err)
	}

	switch levels[0] {
	case "admin":
		return usecases.UserRoleAdmin, nil
	case "head":
		return usecases.UserRoleHead, nil
	case "staff":
		return usecases.UserRoleStaff, nil
	default:
		return usecases.UserRoleNone, errors.New("unknown role")
	}
}

func NewGetUserRoleUseCase(db *sqlx.DB) usecases.GetUserRoleUseCase {
	return &getUserRole{
		db: db,
	}
}
