package notes

import (
	"context"
	"microservices_course/week4/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.Note, error) {
	note, err := s.noteRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return note, nil
}
