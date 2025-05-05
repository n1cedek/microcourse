package notea

import (
	"context"
	"log"
	"microservices_course/week3/internal/converter"
	desc "microservices_course/week3/pkg/note_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	//Делаем запрос на получение данных
	noteO, err := i.noteService.Get(ctx, req.GetId())

	if err != nil {
		return nil, err
	}
	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n",
		noteO.ID, noteO.Info.Title, noteO.Info.Content, noteO.CreatedAt, noteO.UpdatedAt)

	return &desc.GetResponse{Note: converter.ToNoteFromServ(noteO)}, nil

}
