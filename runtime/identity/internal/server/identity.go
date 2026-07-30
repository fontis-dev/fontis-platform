package server

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	pb "github.com/fontis-dev/fontis-platform/runtime/identity/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateHousehold(ctx context.Context, req *pb.CreateHouseholdRequest) (*pb.CreateHouseholdResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	h, err := s.store.CreateHousehold(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create household: %v", err)
	}

	return &pb.CreateHouseholdResponse{Household: h}, nil
}

func (s *Server) GetHousehold(ctx context.Context, req *pb.GetHouseholdRequest) (*pb.GetHouseholdResponse, error) {
	if req.HouseholdId == "" {
		return nil, status.Error(codes.InvalidArgument, "household_id is required")
	}

	h, err := s.store.GetHousehold(ctx, req.HouseholdId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "household %s not found", req.HouseholdId)
		}
		return nil, status.Errorf(codes.Internal, "get household: %v", err)
	}

	return &pb.GetHouseholdResponse{Household: h}, nil
}

func (s *Server) UpdateHousehold(ctx context.Context, req *pb.UpdateHouseholdRequest) (*pb.UpdateHouseholdResponse, error) {
	if req.HouseholdId == "" {
		return nil, status.Error(codes.InvalidArgument, "household_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	h, err := s.store.UpdateHousehold(ctx, req.HouseholdId, req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "household %s not found", req.HouseholdId)
		}
		return nil, status.Errorf(codes.Internal, "update household: %v", err)
	}

	return &pb.UpdateHouseholdResponse{Household: h}, nil
}

func (s *Server) DeleteHousehold(ctx context.Context, req *pb.DeleteHouseholdRequest) (*pb.DeleteHouseholdResponse, error) {
	if req.HouseholdId == "" {
		return nil, status.Error(codes.InvalidArgument, "household_id is required")
	}

	if err := s.store.DeleteHousehold(ctx, req.HouseholdId); err != nil {
		return nil, status.Errorf(codes.Internal, "delete household: %v", err)
	}

	return &pb.DeleteHouseholdResponse{}, nil
}

func (s *Server) ListHouseholds(ctx context.Context, req *pb.ListHouseholdsRequest) (*pb.ListHouseholdsResponse, error) {
	households, err := s.store.ListHouseholds(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list households: %v", err)
	}

	return &pb.ListHouseholdsResponse{Households: households}, nil
}

func (s *Server) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	if req.HouseholdId == "" {
		return nil, status.Error(codes.InvalidArgument, "household_id is required")
	}
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display_name is required")
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	p, err := s.store.CreateProfile(ctx, req.HouseholdId, req.DisplayName, role)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return nil, status.Errorf(codes.NotFound, "household %s not found", req.HouseholdId)
		}
		return nil, status.Errorf(codes.Internal, "create profile: %v", err)
	}

	return &pb.CreateProfileResponse{Profile: p}, nil
}

func (s *Server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if req.ProfileId == "" {
		return nil, status.Error(codes.InvalidArgument, "profile_id is required")
	}

	p, err := s.store.GetProfile(ctx, req.ProfileId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "profile %s not found", req.ProfileId)
		}
		return nil, status.Errorf(codes.Internal, "get profile: %v", err)
	}

	return &pb.GetProfileResponse{Profile: p}, nil
}

func (s *Server) ListProfiles(ctx context.Context, req *pb.ListProfilesRequest) (*pb.ListProfilesResponse, error) {
	if req.HouseholdId == "" {
		return nil, status.Error(codes.InvalidArgument, "household_id is required")
	}

	profiles, err := s.store.ListProfiles(ctx, req.HouseholdId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list profiles: %v", err)
	}

	return &pb.ListProfilesResponse{Profiles: profiles}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if req.ProfileId == "" {
		return nil, status.Error(codes.InvalidArgument, "profile_id is required")
	}

	displayName := req.DisplayName
	role := req.Role

	p, err := s.store.UpdateProfile(ctx, req.ProfileId, displayName, role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "profile %s not found", req.ProfileId)
		}
		return nil, status.Errorf(codes.Internal, "update profile: %v", err)
	}

	return &pb.UpdateProfileResponse{Profile: p}, nil
}

func (s *Server) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error) {
	if req.ProfileId == "" {
		return nil, status.Error(codes.InvalidArgument, "profile_id is required")
	}

	if err := s.store.DeleteProfile(ctx, req.ProfileId); err != nil {
		return nil, status.Errorf(codes.Internal, "delete profile: %v", err)
	}

	return &pb.DeleteProfileResponse{}, nil
}
