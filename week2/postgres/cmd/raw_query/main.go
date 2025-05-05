package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"log"
	"microservices_course/postgres/internal/config"
	"microservices_course/postgres/internal/config/env"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbConfig, err := env.NewDbConfig()
	if err != nil {
		log.Fatalf("failed to get db config: %v", err)
	}
	//СОЗДАЕМ СОЕДИНЕНИЕ С БД
	con, err := pgx.Connect(ctx, dbConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	//Делаем запрос на вставку
	res, err := con.Exec(ctx, "INSERT INTO note (title,body) VALUES ($1,$2)", gofakeit.BeerName(), gofakeit.Name())
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}
	log.Printf("inserted %d rows", res.RowsAffected())

	//Делаем запрос на выборку записей ищ таблицы
	rows, err := con.Query(ctx, "SELECT id, title, body, created_at, updated_at FROM note")
	if err != nil {
		log.Fatalf("failed to selected rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, body string
		var createdAt time.Time
		var updateAt sql.NullTime

		err = rows.Scan(&id, &title, &body, &createdAt, &updateAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}
		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updateAt)
	}
}
