package user_service

import (
	"context"

	"github.com/nandanurseptama/golang-microservices/internal/grpc"
	"github.com/nandanurseptama/golang-microservices/pkg/firestore"
	pb "github.com/nandanurseptama/golang-microservices/services/user_service/api/proto"
	"github.com/nandanurseptama/golang-microservices/services/user_service/dtos"
	"github.com/nandanurseptama/golang-microservices/services/user_service/interfaces"
	"github.com/nandanurseptama/golang-microservices/services/user_service/internal"
	"github.com/sirupsen/logrus"
	grpcLib "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceConfig struct {
	Environment    string
	ServiceId      string
	Address        string
	PasswordSecret string
	PasswordIv     string
}
type service struct {
	pb.UnimplementedUserServiceServer
	usecases interfaces.Usecases
	config   UserServiceConfig
	logger   *logrus.Entry
}

func NewService(
	config UserServiceConfig,
	firestoreClient firestore.FirestoreClient,
) *service {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})
	ctxLogger := logger.WithFields(logrus.Fields{
		"service": config.ServiceId,
	})
	return &service{
		usecases: internal.NewUsecases(
			firestoreClient,
			ctxLogger,
			config.PasswordSecret,
			config.PasswordIv,
		),
		config: config,
		logger: ctxLogger,
	}
}

func (svc *service) RegisterGrpcServer(server *grpcLib.Server) {
	pb.RegisterUserServiceServer(server, svc)
}

func (svc *service) ListenForConnections(ctx context.Context) {
	grpc.ListenForConnections(ctx, svc, svc.config.Address, svc.config.ServiceId)
}

func (svc *service) GetUserByEmail(
	ctx context.Context,
	request *pb.GetUserByEmailRequest,
) (*pb.User, error) {

	user, err := svc.usecases.GetUserByEmail(ctx, request.GetEmail())

	if err != nil {
		svc.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Failed to get user by email")

		return nil, status.Error(codes.Internal, "Failed to get user by email")
	}

	if user == nil {
		svc.logger.Info("User not found")
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return &pb.User{
		Id:        user.Id,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}, err
}

func (svc *service) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.User, error) {
	result, err := svc.usecases.CreateUser(ctx, dtos.CreateUserDto{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, err.ToGrpcError()
	}

	return &pb.User{
		Id:        result.Id,
		Email:     result.Email,
		Password:  result.Password,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		DeletedAt: result.DeletedAt,
	}, nil
}
