package repository

import (
	"context"
	"microservices_course/week4/internal/model"
)

type NoteRepo interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}
