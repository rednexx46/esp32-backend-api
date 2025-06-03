package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var DBName string

func InitDB() {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	DBName = os.Getenv("MONGO_DATABASE")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPass, mongoHost, mongoPort)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("[MongoDB] Failed to connect: %v", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("[MongoDB] Ping failed: %v", err)
	}

	MongoClient = client
	log.Println("[MongoDB] Connected successfully.")
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}
