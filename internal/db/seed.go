package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/rednexx46/esp32-backend-api/internal/models"
	"github.com/rednexx46/esp32-backend-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func SeedAdminUser() {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		log.Println("[SEED] ADMIN_USERNAME or ADMIN_PASSWORD not set.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users := MongoClient.Database(DBName).Collection("users")

	count, err := users.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		log.Printf("[SEED] Failed to count users: %v", err)
		return
	}

	if count > 0 {
		log.Println("[SEED] Admin user already exists. Skipping seed.")
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("[SEED] Failed to hash password: %v", err)
		return
	}

	admin := models.User{
		Username: username,
		Password: hashedPassword,
		Role:     "admin",
	}

	_, err = users.InsertOne(ctx, admin)
	if err != nil {
		log.Printf("[SEED] Failed to insert admin user: %v", err)
		return
	}

	log.Println("[SEED] Admin user created successfully.")
}
