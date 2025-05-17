package main

import (
	"context"
	"flag"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"microservices_course/config/internal/config"
	"microservices_course/config/internal/config/env"
	"microservices_course/week7/grpcl/grpc/pkg/note_v1"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	note_v1.UnimplementedNoteV1Server
	pool *pgxpool.Pool
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
	flag.Parse()
	//ctx := context.Background()

	//Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	//pgConfig, err := env.NewDsnConfig()
	//if err != nil {
	//	log.Fatalf("failed to get pg config: %v", err)
	//}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen server: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	note_v1.RegisterNoteV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
