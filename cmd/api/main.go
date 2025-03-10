package main

import (
	"context"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"os"
	"time"

	"BE-HA001/pkg/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := storage.NewStorage()
	if err != nil {
		logger.Log(logger.FatalLevel, "Failed to connect to database")
	}

	httpServer := NewHTTPServer(port, db)

	// Tạo Go Micro service
	service := micro.NewService(
		micro.Name("go.micro.httpserver"),
		micro.Version("latest"),
		micro.BeforeStart(func() error {
			httpServer.Start()
			return nil
		}),
		micro.BeforeStop(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			httpServer.Stop(ctx)
			return nil
		}),
	)

	if err := service.Run(); err != nil {
		logger.Log(logger.FatalLevel, "Service failed")
	}
}
