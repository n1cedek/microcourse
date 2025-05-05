package note

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"microservices_course/week3/internal/model"
	"microservices_course/week3/internal/repo"
	"microservices_course/week3/internal/repo/note/converter"
	modelRepo "microservices_course/week3/internal/repo/note/model"
)

const (
	tableName = "note"

	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) repository.NoteRepo {
	return &repo{db: db}

}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	//Делаем запрос на вставку записей в таблицу
	builderInsert := sq.Insert(tableName).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to builed query: %v", err)
	}
	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}
	log.Printf("inserted note with noteID: %d", id)

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	//Делаем запрос на получение данных
	builderSelect := sq.Select(idColumn, titleColumn, contentColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		OrderBy("id ASC").
		Limit(10)

	query, arg, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var note modelRepo.Note
	err = r.db.QueryRow(ctx, query, arg...).Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil

}
