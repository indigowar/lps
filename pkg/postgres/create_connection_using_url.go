package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateConnectionUsingURL(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}
