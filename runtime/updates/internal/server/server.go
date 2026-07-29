package server

import (
	"google.golang.org/grpc"

	"github.com/fontis-dev/fontis-platform/runtime/updates/internal/config"
)

type Server struct {
	grpc *grpc.Server
	cfg  *config.Config
}

func New(cfg *config.Config) *grpc.Server {
	srv := grpc.NewServer()
	_ = &Server{cfg: cfg, grpc: srv}
	return srv
}
