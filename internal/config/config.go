package config

import (
	"os"
)

type Config struct {
	Db struct {
		SystemUser     string
		SystemPassword string

		AdminUser     string
		AdminPassword string

		HeadUser     string
		HeadPassword string

		StaffUser     string
		StaffPassword string

		Host string
		Db   string
		Port string
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.Db.Host = os.Getenv("DB_HOST")
	cfg.Db.Db = os.Getenv("DB_NAME")
	cfg.Db.Port = os.Getenv("DB_PORT")

	cfg.Db.SystemUser = os.Getenv("POSTGRES_SYSTEM_NAME")
	cfg.Db.SystemPassword = os.Getenv("POSTGRES_SYSTEM_PWD")

	cfg.Db.AdminUser = os.Getenv("DB_ADMIN_NAME")
	cfg.Db.AdminPassword = os.Getenv("DB_ADMIN_PWD")

	cfg.Db.HeadUser = os.Getenv("DB_HEAD_NAME")
	cfg.Db.HeadPassword = os.Getenv("DB_HEAD_PWD")

	cfg.Db.StaffUser = os.Getenv("DB_STAFF_NAME")
	cfg.Db.StaffPassword = os.Getenv("DB_STAFF_PWD")

	return cfg, nil
}
