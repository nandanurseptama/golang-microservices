package auth_service

import (
	"github.com/nandanurseptama/golang-microservices/internal/grpc"
	auth_service_proto "github.com/nandanurseptama/golang-microservices/services/auth_service/api/proto"
)

// Create auth service client
func NewClient(address string) auth_service_proto.AuthServiceClient {
	conn := grpc.CreateClientConnection(address)
	return auth_service_proto.NewAuthServiceClient(conn)
}
