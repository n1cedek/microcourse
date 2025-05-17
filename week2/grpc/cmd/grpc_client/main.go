package main

import (
	"context"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"microservices_course/grpc/pkg/note_v1"
	"time"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := note_v1.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &note_v1.GetRequest{Id: noteID})

	if err != nil {
		log.Fatalf("failed to get note by ID: %v", err)
	}

	log.Printf(color.RedString("Note info:\n", color.GreenString("%+v", r.GetNote())))

}
