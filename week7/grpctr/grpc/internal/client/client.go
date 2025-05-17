package client

import (
	"context"
	"microservices_course/week7/grpctr/grpc/internal/model"
)

type OtherServiceClient interface {
	Get(ctx context.Context, id int64) (*model.Note, error)
}
