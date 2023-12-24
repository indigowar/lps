package login

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Login(ctx context.Context, string, password string) (uuid.UUID, error)
}
