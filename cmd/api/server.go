package main

import (
	router "BE-HA001/cmd/api/route"
	"BE-HA001/pkg/export"
	"BE-HA001/pkg/storage"
	"context"
	"go-micro.dev/v4/logger"
	"net/http"
	"time"
)

type HTTPServer struct {
	httpServer *http.Server
	db         *storage.Storage
}

func NewHTTPServer(port string, db *storage.Storage) *HTTPServer {
	defaultExport := &export.PDF{}
	return &HTTPServer{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      router.NewRouter(db, defaultExport),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		db: db,
	}
}

func (s *HTTPServer) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log(logger.FatalLevel, "Server error")
		}
	}()
}

func (s *HTTPServer) Stop(ctx context.Context) {
	logger.Log(logger.InfoLevel, "Shutting down HTTP server...")

	if s.db != nil {
		s.db.Close()
		logger.Log(logger.InfoLevel, "Database connection closed")
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Log(logger.ErrorLevel, "Error shutting down HTTP server")
	} else {
		logger.Log(logger.InfoLevel, "HTTP server stopped gracefully")
	}
}
