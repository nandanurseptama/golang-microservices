package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nandanurseptama/golang-microservices/internal/failure"
	"github.com/nandanurseptama/golang-microservices/pkg/crypto"
	"github.com/nandanurseptama/golang-microservices/pkg/firestore"
	"github.com/nandanurseptama/golang-microservices/pkg/validator"
	"github.com/nandanurseptama/golang-microservices/services/user_service/dtos"
	"github.com/nandanurseptama/golang-microservices/services/user_service/entities"
	"github.com/nandanurseptama/golang-microservices/services/user_service/interfaces"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type UsecasesImpl struct {
	database   firestore.FirestoreClient
	logger     *logrus.Entry
	validator  *validator.Validator
	aes        crypto.Aes
	passwordIv string
}

func NewUsecases(
	database firestore.FirestoreClient,
	logger *logrus.Entry,
	passwordSecret string,
	passwordIv string,
) interfaces.Usecases {
	v := validator.New()
	aes := crypto.NewAes(passwordSecret)
	return &UsecasesImpl{
		database:   database,
		logger:     logger,
		validator:  v,
		aes:        aes,
		passwordIv: passwordIv,
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
func (u *UsecasesImpl) CreateUser(
	ctx context.Context,
	dto dtos.CreateUserDto,
) (*entities.UserEntity, *failure.Failure) {
	err := u.validator.ValidateStruct(dto)

	if err != nil {

		return nil, failure.New(
			failure.ValidationError,
			err.Error(),
		)
	}
	dto.Email = strings.ToLower(dto.Email)
	encryptedPassword, err := u.aes.EncryptCBC(dto.Password, u.passwordIv)

	if err != nil {
		u.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(
			fmt.Sprintf("Failed to encrypt password. Plain password : '%s'", dto.Password),
		)

		return nil, failure.New(
			failure.ValidationError,
			"Internal error",
		)
	}
	if encryptedPassword == nil {
		u.logger.WithFields(logrus.Fields{
			"error": "Encrypted password return nil",
		}).Error(
			fmt.Sprintf("Failed to encrypt password. Plain password : '%s'", dto.Password),
		)
		return nil, failure.New(
			failure.ValidationError,
			"Internal error",
		)
	}
	if len(*encryptedPassword) < 1 {
		u.logger.WithFields(logrus.Fields{
			"error": "Encrypted password return empty",
		}).Error(
			fmt.Sprintf("Failed to encrypt password. Plain password : '%s'", dto.Password),
		)
		return nil, failure.New(
			failure.ValidationError,
			"Internal error",
		)
	}

	dto.Password = *encryptedPassword

	result, err := u.database.GetCollection("users").WhereColumn(
		"email",
		firestore.Operator.Equal,
		dto.Email,
	).GetOne(ctx)

	if err != nil {

		u.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to check email on database")

		return nil, failure.New(
			failure.InternalError,
			err.Error(),
		)
	}

	if result != nil {
		u.logger.WithFields(logrus.Fields{
			"error": fmt.Sprintf("User with email '%s' already registered", dto.Email),
		}).Error(
			fmt.Sprintf("User with email '%s' already registered", dto.Email),
		)

		return nil, failure.New(
			failure.ValidationError,
			fmt.Sprintf("User with email '%s' already registered", dto.Email),
		)
	}

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
