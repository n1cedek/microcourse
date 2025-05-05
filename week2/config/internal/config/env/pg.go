package env

import (
	"errors"
	"microservices_course/config/internal/config"
	"os"
)

var _ config.PGConfig = (*pgConfig)(nil)

const pgDsnEnv = "PG_DSN"

type pgConfig struct {
	dsn string
}

func NewDsnConfig() (*pgConfig, error) {
	dsn := os.Getenv(pgDsnEnv)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}
	return &pgConfig{dsn: dsn}, nil
}

func (d *pgConfig) DSN() string {
	return d.dsn
}
