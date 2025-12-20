// Package grpc
package grpc

import (
	"user-service/src/interfaces/content"
	pb "user-service/src/interfaces/grpc"

	"half-nothing.cn/service-core/interfaces/logger"
)

type AuthServer struct {
	pb.UnimplementedAuthServer
	logger logger.Interface
}

func NewAuthServer(
	content *content.ApplicationContent,
) *AuthServer {
	return &AuthServer{
		logger: logger.NewLoggerAdapter(content.Logger(), "grpc-server"),
	}
}
