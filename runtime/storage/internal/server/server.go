package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/fontis-dev/fontis-platform/runtime/storage/internal/config"
	pb "github.com/fontis-dev/fontis-platform/runtime/storage/proto"
)

type Server struct {
	pb.UnimplementedStorageServiceServer
	cfg  *config.Config
	grpc *grpc.Server
}

func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Start(ctx context.Context) error {
	if err := os.Remove(s.cfg.SocketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove socket: %w", err)
	}

	lis, err := net.Listen("unix", s.cfg.SocketPath)
	if err != nil {
		return fmt.Errorf("listen on unix socket: %w", err)
	}

	s.grpc = grpc.NewServer()
	pb.RegisterStorageServiceServer(s.grpc, s)
	reflection.Register(s.grpc)

	log.Printf("storage service listening on %s", s.cfg.SocketPath)

	go func() {
		if err := s.grpc.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.grpc.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		return nil
	case <-ctx.Done():
		s.grpc.Stop()
		return ctx.Err()
	}
}
