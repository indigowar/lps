package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateConnection(host string, port string, db string, user string, password string) (*sqlx.DB, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)
	return sqlx.Connect("postgres", url)
}
