package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var sensorsCollection *mongo.Collection

func decryptPayloadIfNeeded(payload string) string {
	if os.Getenv("ENCRYPTION") != "true" {
		return payload
	}

	cipherURL := os.Getenv("ENCRYPT_API_URL")
	reqBody, _ := json.Marshal(map[string]string{
		"payload": payload,
	})

	resp, err := http.Post(cipherURL+"decrypt", "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		return "decryption_failed"
	}
	defer resp.Body.Close()

	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "decryption_error"
	}
	return res["decrypted"]
}

// GetSensorDataByDevice godoc
// @Summary      Get sensor data by device ID
// @Description  Retrieves all sensor data associated with a specific device ID. Decrypts the payload if needed before returning.
// @Tags         sensors
// @Param        device_id  path      string  true  "Device ID"
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /sensors/{device_id} [get]
func GetSensorDataByDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	filter := bson.M{"device_id": deviceID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := sensorsCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding data"})
		return
	}

	for _, doc := range results {
		if payload, ok := doc["payload"].(string); ok {
			doc["payload"] = decryptPayloadIfNeeded(payload)
		}
	}

	c.JSON(http.StatusOK, results)
}

// GetAllSensorData godoc
// @Summary      Retrieve all sensor data
// @Description  Fetches all sensor data from the database, decrypting payloads if necessary.
// @Tags         sensors
// @Produce      json
// @Success      200  {array}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /sensors [get]
func GetAllSensorData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := sensorsCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding data"})
		return
	}

	for _, doc := range results {
		if payload, ok := doc["payload"].(string); ok {
			doc["payload"] = decryptPayloadIfNeeded(payload)
		}
	}

	c.JSON(http.StatusOK, results)
}

// GetActiveDevices godoc
// @Summary      Get active devices
// @Description  Retrieves a list of unique active device IDs from the sensors collection.
// @Tags         devices
// @Produce      json
// @Success      200  {array}  map[string]interface{}  "List of active device IDs"
// @Failure      500  {object}  map[string]string      "Internal server error"
// @Router       /devices/active [get]
func GetActiveDevices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id": "$device_id",
		}}},
	}
	cursor, err := sensorsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Aggregation failed"})
		return
	}
	defer cursor.Close(ctx)

	var devices []bson.M
	if err := cursor.All(ctx, &devices); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding devices"})
		return
	}

	c.JSON(http.StatusOK, devices)
}
