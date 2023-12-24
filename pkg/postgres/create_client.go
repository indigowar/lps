package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Credentials struct {
	Host     string
	Port     string
	Db       string
	User     string
	Password string
}

func CreateClient(c Credentials) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Db)
	return pgx.Connect(context.Background(), url)
}
