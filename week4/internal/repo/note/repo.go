package note

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"log"
	"microservices_course/week4/internal/client/db"
	"microservices_course/week4/internal/model"
	"microservices_course/week4/internal/repo"
	"microservices_course/week4/internal/repo/note/converter"
	modelRepo "microservices_course/week4/internal/repo/note/model"
)

const (
	tableName       = "note"
	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "body"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepo(db db.Client) repository.NoteRepo {
	return &repo{db: db}

}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	//Делаем запрос на вставку записей в таблицу
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to builed query: %v", err)
	}

	q := db.Query{
		Name:     "note_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	//Делаем запрос на получение данных
	builderSelect := sq.Select(idColumn, titleColumn, contentColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, arg, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "note_repository.Get",
		QueryRaw: query,
	}

	var note modelRepo.Note
	err = r.db.DB().ScanOneContext(ctx, &note, q, arg...)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil

}
