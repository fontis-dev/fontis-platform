package server

import (
	"context"

	pb "github.com/fontis-dev/fontis-platform/runtime/networking/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListInterfaces(ctx context.Context, req *pb.ListInterfacesRequest) (*pb.ListInterfacesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListInterfaces not implemented")
}

func (s *Server) ScanWiFiNetworks(ctx context.Context, req *pb.ScanWiFiNetworksRequest) (*pb.ScanWiFiNetworksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ScanWiFiNetworks not implemented")
}

func (s *Server) ConnectWiFi(ctx context.Context, req *pb.ConnectWiFiRequest) (*pb.ConnectWiFiResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ConnectWiFi not implemented")
}

func (s *Server) DisconnectWiFi(ctx context.Context, req *pb.DisconnectWiFiRequest) (*pb.DisconnectWiFiResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "DisconnectWiFi not implemented")
}

func (s *Server) GetConnectionStatus(ctx context.Context, req *pb.GetConnectionStatusRequest) (*pb.GetConnectionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetConnectionStatus not implemented")
}
