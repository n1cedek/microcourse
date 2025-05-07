package notes

import (
	"context"
	"microservices_course/week3/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.noteRepo.Create(ctx, info)
		if errTx != nil {
			return errTx
		}
		_, errTx = s.noteRepo.Get(ctx, id)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
