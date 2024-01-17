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
	repo           repository.RepositoryInterface
	tokenGenerator auth.TokenGenerator
}

func NewGrpcServer(repo repository.RepositoryInterface, tokenGenerator auth.TokenGenerator) *GrpcServer {
	return &GrpcServer{repo: repo, tokenGenerator: tokenGenerator}
}

func (s *GrpcServer) IsAuthenticated(ctx context.Context, r *proto.AuthRequest) (*proto.AuthResponse, error) {
	tokenString := r.GetToken()

	claims, err := s.tokenGenerator.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	exp := int64(claims["exp"].(float64))
	expireTime := time.Unix(exp, 0)

	if time.Now().After(expireTime) {
		return &proto.AuthResponse{IsAuthenticated: false}, nil
	}

	return &proto.AuthResponse{IsAuthenticated: true}, nil

}
