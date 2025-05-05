package converter

import (
	modelRepo "microservices_course/week3/internal/repo/note/model"

	"microservices_course/week3/internal/model"
)

func ToNoteFromRepo(note *modelRepo.Note) *model.Note {
	return &model.Note{
		ID:        note.ID,
		Info:      ToNoteInfoRepo(note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToNoteInfoRepo(info modelRepo.NoteInfo) model.NoteInfo {
	return model.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
