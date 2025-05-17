package other_service

import (
	"context"
	"database/sql"
	"microservices_course/week7/grpctr/grpc/internal/model"
	desc "microservices_course/week7/grpctr/grpc/pkg/other_note_v1"
)

type client struct {
	noteClient desc.OtherNoteV1Client
}

func New(noteClient desc.OtherNoteV1Client) *client {
	return &client{noteClient: noteClient}
}

func (c *client) Get(ctx context.Context, id int64) (*model.Note, error) {
	res, err := c.noteClient.Get(ctx, &desc.GetRequest{Id: id})
	if err != nil {
		return nil, err
	}

	var updatedAt sql.NullTime
	if res.GetNote().UpdatedAt != nil {
		updatedAt = sql.NullTime{
			Time:  res.GetNote().GetUpdatedAt().AsTime(),
			Valid: true,
		}
	}

	return &model.Note{
		ID: res.GetNote().GetId(),
		Info: model.NoteInfo{
			Title:   res.GetNote().GetInfo().GetTitle(),
			Content: res.GetNote().GetInfo().GetContent(),
		},
		CreatedAt: res.GetNote().GetCreatedAt().AsTime(),
		UpdatedAt: updatedAt,
	}, nil
}
