package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/config"
	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/store"
	"github.com/fontis-dev/fontis-platform/runtime/identity/proto"
	fontistls "github.com/fontis-dev/fontis-platform/runtime/pkg/tls"
)

type Server struct {
	proto.UnimplementedIdentityServiceServer
	cfg   *config.Config
	store *store.Store
	grpc  *grpc.Server
}

func New(cfg *config.Config, st *store.Store, tlsCfg *tls.Config) *Server {
	opts := []grpc.ServerOption{}
	if tlsCfg != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsCfg)))
	}
	return &Server{
		cfg:   cfg,
		store: st,
		grpc:  grpc.NewServer(opts...),
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.store.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate store: %w", err)
	}

	if err := os.Remove(s.cfg.SocketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove socket: %w", err)
	}

	lis, err := net.Listen("unix", s.cfg.SocketPath)
	if err != nil {
		return fmt.Errorf("listen on unix socket: %w", err)
	}

	proto.RegisterIdentityServiceServer(s.grpc, s)
	reflection.Register(s.grpc)

	log.Printf("identity service listening on %s", s.cfg.SocketPath)

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

func LoadTLSConfig(cfg *config.Config) (*tls.Config, error) {
	if cfg.DevTLS {
		serverCfg, _, err := fontistls.DevTLSConfig("identity")
		if err != nil {
			return nil, fmt.Errorf("dev tls config: %w", err)
		}
		return serverCfg, nil
	}

	caPEM, err := os.ReadFile(cfg.TLSCAPath)
	if err != nil {
		return nil, fmt.Errorf("read ca cert: %w", err)
	}
	certPEM, err := os.ReadFile(cfg.TLSCertPath)
	if err != nil {
		return nil, fmt.Errorf("read cert: %w", err)
	}
	keyPEM, err := os.ReadFile(cfg.TLSKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read key: %w", err)
	}

	tlsCfg, err := fontistls.NewServerTLSConfig(caPEM, certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("server tls config: %w", err)
	}
	return tlsCfg, nil
}


