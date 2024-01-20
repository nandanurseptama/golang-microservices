package internal

import (
	"context"
	"time"

	"github.com/nandanurseptama/golang-microservices/internal/failure"
	"github.com/nandanurseptama/golang-microservices/pkg/firestore"
	"github.com/nandanurseptama/golang-microservices/services/user_service/dtos"
	"github.com/nandanurseptama/golang-microservices/services/user_service/entities"
	"github.com/nandanurseptama/golang-microservices/services/user_service/interfaces"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type UsecasesImpl struct {
	database firestore.FirestoreClient
	logger   *logrus.Entry
}

func NewUsecases(database firestore.FirestoreClient, logger *logrus.Entry) interfaces.Usecases {
	return &UsecasesImpl{
		database: database,
		logger:   logger,
	}
}

func (u *UsecasesImpl) GetUserByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	iter := u.database.
		Collection("users").
		Where("email", "==", email).Limit(1).Documents(ctx)
	var user *entities.UserEntity
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		logrus.Info(doc.Data())
		user, err = entities.UserEntityFromMap(doc.Data())
		user.Id = doc.Ref.ID

		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
func (u *UsecasesImpl) CreateUser(ctx context.Context, dto dtos.CreateUserDto) (*entities.UserEntity, *failure.Failure) {
	creatUserData, err := dto.ToMap()

	if err != nil {

		u.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to convert dto to map")

		return nil, failure.New(
			failure.InternalError,
			"Failed to create new data",
		)
	}

	creatUserData["createdAt"] = time.Now().UTC().String()
	creatUserData["updatedAt"] = nil
	creatUserData["deletedAt"] = nil

	ref, _, err := u.database.Collection("users").
		Add(ctx, creatUserData)

	if err != nil {

		u.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create new data")

		return nil, failure.New(
			failure.InternalError,
			"Failed to create new data",
		)
	}
	doc, err := ref.Get(ctx)

	if err != nil {

		u.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to get new data")

		return nil, failure.New(
			failure.InternalError,
			"Failed to create new data",
		)
	}

	user, err := entities.UserEntityFromMap(doc.Data())

	if err != nil {

		u.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to parsing new data")

		return nil, failure.New(
			failure.InternalError,
			"Failed to create new data",
		)
	}

	user.Id = doc.Ref.ID
	return user, nil
}
