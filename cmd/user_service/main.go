package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/nandanurseptama/golang-microservices/pkg/firestore"
	"github.com/nandanurseptama/golang-microservices/services/user_service"
	"github.com/sirupsen/logrus"
)

const (
	basePath string = "../../"
	envPath         = basePath + ".env"
)

func main() {
	logrus.New()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
	err := godotenv.Load(envPath)
	if err != nil {
		logrus.Error("failed to load env", err)
		return
	}
	firestoreFileName := os.Getenv("firebase_service_account_path")

	firestoreClient, err := firestore.NewClient(basePath + firestoreFileName)

	if err != nil {
		logrus.Error("failed to init firestore client", err)
		return
	}
	serviceId := os.Getenv("user.service.id")
	address := os.Getenv("user.service.address")
	service := user_service.NewService(address, serviceId, *firestoreClient)

	service.ListenForConnections(context.TODO())

}
