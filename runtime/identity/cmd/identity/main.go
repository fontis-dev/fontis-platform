package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "modernc.org/sqlite"

	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/config"
	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/server"
	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/store"
)

func main() {
	cfg := config.Load()

	tlsCfg, err := server.LoadTLSConfig(cfg)
	if err != nil {
		log.Fatalf("load tls config: %v", err)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	db.SetMaxOpenConns(1)

	st := store.New(db)
	srv := server.New(cfg, st, tlsCfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("start server: %v", err)
	}

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	if err := st.Close(); err != nil {
		log.Fatalf("close store: %v", err)
	}
	log.Println("shutdown complete")
	os.Exit(0)
}
