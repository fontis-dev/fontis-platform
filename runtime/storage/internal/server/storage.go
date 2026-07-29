package server

import (
	"context"

	pb "github.com/fontis-dev/fontis-platform/runtime/storage/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListDevices(ctx context.Context, req *pb.ListDevicesRequest) (*pb.ListDevicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListDevices not implemented")
}

func (s *Server) CreateVolume(ctx context.Context, req *pb.CreateVolumeRequest) (*pb.CreateVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "CreateVolume not implemented")
}

func (s *Server) GetVolume(ctx context.Context, req *pb.GetVolumeRequest) (*pb.GetVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetVolume not implemented")
}

func (s *Server) MountVolume(ctx context.Context, req *pb.MountVolumeRequest) (*pb.MountVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "MountVolume not implemented")
}

func (s *Server) UnmountVolume(ctx context.Context, req *pb.UnmountVolumeRequest) (*pb.UnmountVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "UnmountVolume not implemented")
}

func (s *Server) ListVolumes(ctx context.Context, req *pb.ListVolumesRequest) (*pb.ListVolumesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListVolumes not implemented")
}

func (s *Server) DeleteVolume(ctx context.Context, req *pb.DeleteVolumeRequest) (*pb.DeleteVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "DeleteVolume not implemented")
}
