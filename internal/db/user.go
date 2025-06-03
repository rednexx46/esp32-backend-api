package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/rednexx46/esp32-backend-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(user models.User) error {
	collection := GetMongoClient().
		Database(os.Getenv("MONGO_DATABASE")).
		Collection(os.Getenv("MONGO_USERS_COLLECTION"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("[MongoDB] Failed to insert user: %v", err)
		return err
	}
	log.Println("[MongoDB] User created successfully.")
	return nil
}

func FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	DB_NAME := os.Getenv("MONGO_DATABASE")
	if DB_NAME == "" {
		return nil, errors.New("MONGO_DATABASE environment variable is not set")
	}
	USERS_COLLECTION := os.Getenv("MONGO_USERS_COLLECTION")
	if USERS_COLLECTION == "" {
		return nil, errors.New("MONGO_USERS_COLLECTION environment variable is not set")
	}
	collection := MongoClient.Database(DB_NAME).Collection(USERS_COLLECTION)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
