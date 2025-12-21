// Package grpc
package grpc

import (
	pb "user-service/src/interfaces/grpc"

	"half-nothing.cn/service-core/interfaces/logger"
)

type AuthServer struct {
	pb.UnimplementedAuthServer
	logger logger.Interface
}

func NewAuthServer(
	lg logger.Interface,
) *AuthServer {
	return &AuthServer{
		logger: logger.NewLoggerAdapter(lg, "grpc-server"),
	}
}
