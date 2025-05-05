package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"microservices_course/postgres/internal/config"
	"microservices_course/postgres/internal/config/env"
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

	configDb, err := env.NewDbConfig()
	if err != nil {
		log.Fatalf("failed to db config: %v", err)
	}

	//create connecting to database
	con, err := pgx.Connect(ctx, configDb.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer con.Close(ctx)

	//запрос на вставку
	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.BeerName(), gofakeit.Name()).
		Suffix("RETURNING id")

	query, arg, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to builed query: %v", err)
	}

	var noteID int
	err = con.QueryRow(ctx, query, arg...).Scan(&noteID)
	if err != nil {
		log.Fatalf("failed to inserted note: %v", err)
	}
	log.Printf("inserted note with id: %d", noteID)

	//запрос на выборку
	builderSelect := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, arg, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to builed query: %v", err)
	}

	rows, err := con.Query(ctx, query, arg...)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	var id int
	var title, body string
	var createdAt time.Time
	var updateAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &title, &body, &createdAt, &updateAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}
		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n",
			id, title, body, createdAt, updateAt)

	}

	//запрос на обнавление
	builderUpdate := sq.Update("note").
		PlaceholderFormat(sq.Dollar).
		Set("title", gofakeit.BeerName()).
		Set("body", gofakeit.BeerName()).
		Where(sq.Eq{"id": noteID})

	query, arg, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to builed query: %v", err)
	}

	res, err := con.Exec(ctx, query, arg...)
	if err != nil {
		log.Fatalf("failed to update notes: %v", err)
	}
	log.Printf("updated %d rows", res.RowsAffected())

	//запрос на выборку
	builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": noteID}).
		Limit(1)

	query, arg, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to builed query: %v", err)
	}

	err = con.QueryRow(ctx, query, arg...).Scan(&id, &title, &body, &createdAt, &updateAt)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}
	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n",
		id, title, body, createdAt, updateAt)
}
