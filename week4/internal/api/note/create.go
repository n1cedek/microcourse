package notea

import (
	"context"
	"log"
	"microservices_course/week4/internal/converter"
	desc "microservices_course/week4/pkg/note_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	//Делаем запрос на вставку записей в таблицу
	id, err := i.noteService.Create(ctx, converter.ToNoteInfoFromDesc(req.GetInfo()))

	if err != nil {
		return nil, err
	}
	log.Printf("inserted note with noteID: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
