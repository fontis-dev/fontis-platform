package server

import (
	"google.golang.org/grpc"

	"github.com/fontis-dev/fontis-platform/runtime/logging/internal/config"
)

type Server struct {
	grpc *grpc.Server
	cfg  *config.Config
}

func New(cfg *config.Config) *grpc.Server {
	s := &Server{cfg: cfg}
	srv := grpc.NewServer()
	_ = s
	return srv
}
