package main

import (
	"context"
	"flag"
	"log"
	"microservices_course/week3/internal/app"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
