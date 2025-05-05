package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	notea "microservices_course/week3/internal/api/note"
	"microservices_course/week3/internal/config/env"
	"microservices_course/week3/internal/repo/note"
	notes "microservices_course/week3/internal/service/note"
	desc "microservices_course/week3/pkg/note_v1"
	"net"
)

func main() {
	ctx := context.Background()

	//считываем переменные окружения
	err := env.Load(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	dbConfig, err := env.NewDsnConfig()
	if err != nil {
		log.Fatalf("failed to get db config: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcConfig.Address()))
	if err != nil {
		log.Fatalf("failed to listen server: %v", err)
	}

	//Соединение с бд
	con, err := pgxpool.New(ctx, dbConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close()

	noteRepo := note.NewRepo(con)
	noteSrv := notes.NewService(noteRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, notea.NewImplementation(noteSrv))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
