package notes

import (
	repository "microservices_course/week3/internal/repo"
	"microservices_course/week3/internal/service"
)

type serv struct {
	noteRepo repository.NoteRepo
}

func NewService(noteRepo repository.NoteRepo) service.NoteService {
	return &serv{
		noteRepo: noteRepo,
	}
}
