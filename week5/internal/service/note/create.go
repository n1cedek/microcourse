package notes

import (
	"context"
	"microservices_course/week5/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	id, err := s.noteRepo.Create(ctx, info)
	if err != nil {
		return 0, err
	}

	return id, nil
}
