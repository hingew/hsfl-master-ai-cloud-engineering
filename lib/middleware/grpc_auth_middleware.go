package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/proto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"google.golang.org/grpc"
)

type GrpcAuthMiddleware struct {
	client proto.UserServiceClient
}

func NewGrpcAuthMiddleware(conn grpc.ClientConnInterface) *GrpcAuthMiddleware {
	client := proto.NewUserServiceClient(conn)
	return &GrpcAuthMiddleware{client}
}

func (m *GrpcAuthMiddleware) AuthMiddleware(w http.ResponseWriter, r *http.Request, next router.Next) {
	bearerToken := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix("Bearer ", bearerToken)

	if !ok {
		http.Error(w, "No token provied", http.StatusUnauthorized)
	}

	request := &proto.AuthRequest{Token: token}

	res, err := m.client.IsAuthenticated(context.Background(), request)

	if err != nil {
		http.Error(w, "Could not verify token", http.StatusUnauthorized)
	}

	if !res.IsAuthenticated {
		w.WriteHeader(http.StatusUnauthorized)
	}

	next(r)
}
