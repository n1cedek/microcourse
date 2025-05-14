package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"microservices_course/week6/jwt/internal/model"
	descAccess "microservices_course/week6/jwt/pkg/access_v1"
)

var accessToken = flag.String("a", "", "access token")

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	flag.Parse()
	ctx := context.Background()

	md := metadata.New(map[string]string{"Authorization": "Bearer " + *accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	//cred, err := credentials.NewClientTLSFromFile("/home/n1cedek/GolandProjects/microcourse/week6/grpc/service.pem", "")
	//if err != nil {
	//	log.Fatalf("Could not process the credentials: %v", err)
	//}
	//
	//conn, err := grpc.Dial(address, grpc.WithTransportCredentials(cred))
	//if err != nil {
	//	log.Fatalf("failed to connect to server: %v", err)
	//}
	//defer conn.Close()
	//
	//c := note_v1.NewNoteV1Client(conn)
	//
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//
	//r, err := c.Get(ctx, &note_v1.GetRequest{Id: noteID})
	//
	//if err != nil {
	//	log.Fatalf("failed to get note by ID: %v", err)
	//}
	//
	//log.Printf(color.RedString("Note info:\n", color.GreenString("%+v", r.GetNote())))
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	cl := descAccess.NewAccessV1Client(conn)

	_, err = cl.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: model.ExamplePath,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Access granted")
}
