package notes

import (
	"microservices_course/week5/internal/client/db"
	repository "microservices_course/week5/internal/repo"
	"microservices_course/week5/internal/service"
)

type serv struct {
	noteRepo repository.NoteRepo
}

func NewService(noteRepo repository.NoteRepo, txManager db.TxManager) service.NoteService {
	return &serv{
		noteRepo: noteRepo,
	}
}
func NewMockService(deps ...interface{}) service.NoteService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.NoteRepo:
			srv.noteRepo = s
		}
	}

	return &srv
}
