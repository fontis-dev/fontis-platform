package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fontis-dev/fontis-platform/runtime/module-lifecycle/internal/config"
	"github.com/fontis-dev/fontis-platform/runtime/module-lifecycle/internal/server"
)

func main() {
	cfg := config.Load()
	lis, err := net.Listen("unix", cfg.SocketPath)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.SocketPath, err)
	}
	defer lis.Close()

	srv := server.New(cfg)
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	fmt.Println("shutting down module-lifecycle service...")
	srv.GracefulStop()
}
