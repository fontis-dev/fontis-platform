package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fontis-dev/fontis-platform/runtime/networking/internal/config"
	"github.com/fontis-dev/fontis-platform/runtime/networking/internal/server"
)

func main() {
	cfg := config.Load()
	srv := server.New(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("shutdown complete")
	os.Exit(0)
}
