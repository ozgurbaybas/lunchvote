package main

import (
	"context"
	"log"

	"github.com/ozgurbaybas/lunchvote/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("bootstrap application: %v", err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
