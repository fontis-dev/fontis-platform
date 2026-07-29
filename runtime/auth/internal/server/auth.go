package server

import (
	"context"
	"database/sql"
	"errors"
	"time"

	pb "github.com/fontis-dev/fontis-platform/runtime/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if req.ProfileId == "" {
		return nil, status.Error(codes.InvalidArgument, "profile_id is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	ok, err := s.store.VerifyPassword(ctx, req.ProfileId, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "verify password: %v", err)
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "invalid credentials")
	}

	sess, token, refreshToken, err := s.store.CreateSession(ctx, req.ProfileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create session: %v", err)
	}

	resp := &pb.AuthenticateResponse{
		Session: &pb.Session{
			Id:           sess.ID,
			ProfileId:    sess.ProfileID,
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresAt:    timestamppb.New(sess.ExpiresAt),
			CreatedAt:    timestamppb.New(sess.CreatedAt),
		},
	}
	return resp, nil
}

func (s *Server) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	if req.ProfileId == "" {
		return nil, status.Error(codes.InvalidArgument, "profile_id is required")
	}

	sess, token, refreshToken, err := s.store.CreateSession(ctx, req.ProfileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create session: %v", err)
	}

	return &pb.CreateSessionResponse{
		Session: &pb.Session{
			Id:           sess.ID,
			ProfileId:    sess.ProfileID,
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresAt:    timestamppb.New(sess.ExpiresAt),
			CreatedAt:    timestamppb.New(sess.CreatedAt),
		},
	}, nil
}

func (s *Server) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	sess, err := s.store.GetSessionByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.ValidateSessionResponse{Valid: false}, nil
		}
		return nil, status.Errorf(codes.Internal, "validate session: %v", err)
	}

	if time.Now().UTC().After(sess.ExpiresAt) {
		return &pb.ValidateSessionResponse{Valid: false}, nil
	}

	return &pb.ValidateSessionResponse{Valid: true, ProfileId: sess.ProfileID}, nil
}

func (s *Server) RevokeSession(ctx context.Context, req *pb.RevokeSessionRequest) (*pb.RevokeSessionResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	if err := s.store.RevokeSession(ctx, req.SessionId); err != nil {
		return nil, status.Errorf(codes.Internal, "revoke session: %v", err)
	}

	return &pb.RevokeSessionResponse{}, nil
}

func (s *Server) RefreshSession(ctx context.Context, req *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	sess, token, refreshToken, err := s.store.RefreshSession(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "refresh session: %v", err)
	}

	return &pb.RefreshSessionResponse{
		Session: &pb.Session{
			Id:           sess.ID,
			ProfileId:    sess.ProfileID,
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresAt:    timestamppb.New(sess.ExpiresAt),
			CreatedAt:    timestamppb.New(sess.CreatedAt),
		},
	}, nil
}

func (s *Server) CreateApiToken(ctx context.Context, req *pb.CreateApiTokenRequest) (*pb.CreateApiTokenResponse, error) {
	if req.ProfileId == "" {
		return nil, status.Error(codes.InvalidArgument, "profile_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	t, rawToken, err := s.store.CreateAPIToken(ctx, req.ProfileId, req.Name, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create api token: %v", err)
	}

	pbToken := &pb.ApiToken{
		Id:        t.ID,
		ProfileId: t.ProfileID,
		Name:      t.Name,
		TokenHash: rawToken,
		CreatedAt: timestamppb.New(t.CreatedAt),
	}
	if t.ExpiresAt != nil {
		pbToken.ExpiresAt = timestamppb.New(*t.ExpiresAt)
	}

	return &pb.CreateApiTokenResponse{Token: pbToken}, nil
}

func (s *Server) ValidateApiToken(ctx context.Context, req *pb.ValidateApiTokenRequest) (*pb.ValidateApiTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	profileID, err := s.store.ValidateAPIToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateApiTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateApiTokenResponse{Valid: true, ProfileId: profileID}, nil
}

func (s *Server) RevokeApiToken(ctx context.Context, req *pb.RevokeApiTokenRequest) (*pb.RevokeApiTokenResponse, error) {
	if req.TokenId == "" {
		return nil, status.Error(codes.InvalidArgument, "token_id is required")
	}

	if err := s.store.RevokeAPIToken(ctx, req.TokenId); err != nil {
		return nil, status.Errorf(codes.Internal, "revoke api token: %v", err)
	}

	return &pb.RevokeApiTokenResponse{}, nil
}

func (s *Server) ListApiTokens(ctx context.Context, req *pb.ListApiTokensRequest) (*pb.ListApiTokensResponse, error) {
	if req.ProfileId == "" {
		return nil, status.Error(codes.InvalidArgument, "profile_id is required")
	}

	tokens, err := s.store.ListAPITokens(ctx, req.ProfileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list api tokens: %v", err)
	}

	return &pb.ListApiTokensResponse{Tokens: tokens}, nil
}
