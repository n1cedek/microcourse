package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservices_course/week4/internal/model"
	desc "microservices_course/week4/pkg/note_v1"
)

func ToNoteFromServ(note *model.Note) *desc.Note {
	var updateTime *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updateTime = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromServ(note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updateTime,
	}
}

func ToNoteInfoFromServ(info model.NoteInfo) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
func ToNoteInfoFromDesc(info *desc.NoteInfo) *model.NoteInfo {
	return &model.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
