package notes

import (
	"microservices_course/week3/internal/client/db"
	repository "microservices_course/week3/internal/repo"
	"microservices_course/week3/internal/service"
)

type serv struct {
	noteRepo  repository.NoteRepo
	txManager db.TxManager
}

func NewService(noteRepo repository.NoteRepo, txManager db.TxManager) service.NoteService {
	return &serv{
		noteRepo:  noteRepo,
		txManager: txManager,
	}
}
