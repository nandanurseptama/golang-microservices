package interfaces

import (
	"context"

	"github.com/nandanurseptama/golang-microservices/internal/failure"
	"github.com/nandanurseptama/golang-microservices/services/user_service/dtos"
	"github.com/nandanurseptama/golang-microservices/services/user_service/entities"
)

// Interfaces for usecases
type Usecases interface {
	// Get user by email
	//
	// return nil if user not found/
	//
	// return error if any error occured
	GetUserByEmail(ctx context.Context, email string) (*entities.UserEntity, error)

	CreateUser(ctx context.Context, dto dtos.CreateUserDto) (*entities.UserEntity, *failure.Failure)
}
