package server

import (
	"google.golang.org/grpc"

	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/config"
	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/store"
)

type Server struct {
	grpc  *grpc.Server
	cfg   *config.Config
	store *store.Store
}

func New(cfg *config.Config) *grpc.Server {
	srv := grpc.NewServer()
	s := &Server{cfg: cfg, grpc: srv}
	_ = s.store
	return srv
}
