package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"microservices_course/week6/tls/grpc/pkg/note_v1"
	"net"
)

const grpcPort = 50051

type server struct {
	note_v1.UnimplementedNoteV1Server
}

func (s *server) Get(ctx context.Context, req *note_v1.GetRequest) (*note_v1.GetResponse, error) {
	log.Printf("Get ID: %v", req.GetId())

	return &note_v1.GetResponse{
		Note: &note_v1.Note{
			Id: req.GetId(),
			Info: &note_v1.NoteInfo{
				Title:    gofakeit.BeerName(),
				Content:  gofakeit.IPv4Address(),
				Author:   gofakeit.Name(),
				IsPublic: gofakeit.Bool(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen server: %v", err)
	}

	cred, err := credentials.NewServerTLSFromFile("/home/n1cedek/GolandProjects/microcourse/week6/grpc/service.pem", "/home/n1cedek/GolandProjects/microcourse/week6/grpc/service.key")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(cred))
	reflection.Register(s)
	note_v1.RegisterNoteV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
