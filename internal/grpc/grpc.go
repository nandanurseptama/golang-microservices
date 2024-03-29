package grpc

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceRegistrar will make sure our services can be registered as grpc servers
type ServiceRegistrar interface {
	RegisterGrpcServer(server *grpc.Server)
}

// CreateClientConnection will create grpc client
func CreateClientConnection(addr string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logrus.Fatal(err)
	}

	return conn
}

func ListenForConnections(ctx context.Context, registrar ServiceRegistrar, addr, serviceName string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatal(err)
	}

	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)

	registrar.RegisterGrpcServer(srv)

	logrus.Infof("%s running at address %s", serviceName, addr)

	go listenForStopped(ctx, srv, serviceName)

	if err = srv.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}

func listenForStopped(ctx context.Context, grpcServer *grpc.Server, serviceName string) {
	defer func() {
		logrus.Infof("%s stopped", serviceName)
	}()
	for {
		select {
		case <-ctx.Done():
			grpcServer.Stop()
			return
		}
	}
}
