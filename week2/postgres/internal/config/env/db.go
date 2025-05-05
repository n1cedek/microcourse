package env

import (
	"errors"
	"microservices_course/postgres/internal/config"
	"os"
)

var _ config.DBConfig = (*dbConfig)(nil)

const pgDsnEnv = "PG_DSN"
type dbConfig struct {
	dsn string
}

func NewDbConfig()(*dbConfig,error){
	dsn:=os.Getenv(pgDsnEnv)
	if len(dsn) == 0{
		return nil, errors.New("db config not found")
	}
	return &dbConfig{dsn: dsn},nil
}

func (d *dbConfig) DSN()string{
	return d.dsn
}