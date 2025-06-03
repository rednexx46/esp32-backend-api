package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var kpisCollection *mongo.Collection

// GetKPIsByDevice godoc
// @Summary      Get KPIs by device
// @Description  Retrieves KPIs for a specific device_id, paginated
// @Tags         kpis
// @Accept       json
// @Produce      json
// @Param        device_id  path     string  true  "Device ID"
// @Param        limit      query    int     false "Max results (default: 100)"
// @Param        page       query    int     false "Page number (default: 1)"
// @Success      200        {array}  map[string]interface{}
// @Failure      500        {object} map[string]string
// @Router       /api/kpis/device/{device_id} [get]
func GetKPIsByDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	limitStr := c.DefaultQuery("limit", "100")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	skip := (page - 1) * limit

	filter := bson.M{"device_id": deviceID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(skip))

	cursor, err := kpisCollection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch KPIs"})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding KPI data"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetAllKPIs godoc
// @Summary      Get all KPIs
// @Description  Retrieves all KPIs with pagination
// @Tags         kpis
// @Accept       json
// @Produce      json
// @Param        limit  query    int     false  "Max results (default: 100)"
// @Param        page   query    int     false  "Page number (default: 1)"
// @Success      200    {array}  map[string]interface{}
// @Failure      500    {object} map[string]string
// @Router       /api/kpis [get]
func GetAllKPIs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	skip := (page - 1) * limit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(skip))

	cursor, err := kpisCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch KPIs"})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding KPI data"})
		return
	}

	c.JSON(http.StatusOK, results)
}
