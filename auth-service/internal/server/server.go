package server

import (
	"auth-service/internal/repository"
	"auth-service/internal/utils"
	"auth-service/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	repo *repository.UserRepository
	proto.UnimplementedAuthServiceServer
}

func NewAuthServer(repo *repository.UserRepository) *AuthServer {
	return &AuthServer{repo: repo}
}

func (s *AuthServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	err := s.repo.CreateUser(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "Registration failed")
	}
	return &proto.RegisterResponse{Success: true}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil || user == nil || !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}
	token, _ := utils.GenerateJWT(user.ID)
	return &proto.LoginResponse{Token: token}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	claims, err := utils.ValidateJWT(req.Token)
	if err != nil {
		return &proto.ValidateTokenResponse{Valid: false}, nil
	}
	return &proto.ValidateTokenResponse{Valid: true, UserId: claims.UserID}, nil
}
