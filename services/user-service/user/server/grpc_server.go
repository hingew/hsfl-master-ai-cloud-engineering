package server

import (
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/proto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/repository"
	"golang.org/x/net/context"
)

type GrpcServer struct {
	proto.UnimplementedUserServiceServer
	repo           repository.Repository
	tokenGenerator auth.TokenGenerator
}

func NewGrpcServer(repo repository.Repository, tokenGenerator auth.TokenGenerator) *GrpcServer {
	return &GrpcServer{repo: repo}
}

func (s *GrpcServer) IsAuthenticated(ctx context.Context, r *proto.AuthRequest) (*proto.AuthResponse, error) {
	tokenString := r.GetToken()

	claims, err := s.tokenGenerator.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	expireTime := time.Unix(claims["exp"].(int64), 0)

	if time.Now().After(expireTime) {
		return &proto.AuthResponse{IsAuthenticated: false}, nil
	}

	return &proto.AuthResponse{IsAuthenticated: true}, nil

}
