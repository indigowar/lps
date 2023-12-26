package usecases

import (
	"context"

	"github.com/google/uuid"
)

type UserRole int

const (
	UserRoleAdmin UserRole = iota
	UserRoleHead
	UserRoleStaff
	UserRoleNone
)

type GetUserRoleUseCase interface {
	GetUserRole(ctx context.Context, id uuid.UUID) (UserRole, error)
}
