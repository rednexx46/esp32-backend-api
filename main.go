package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/rednexx46/esp32-backend-api/docs"
	"github.com/rednexx46/esp32-backend-api/internal/db"
	"github.com/rednexx46/esp32-backend-api/internal/handlers"
	"github.com/rednexx46/esp32-backend-api/internal/middleware"
	"github.com/rednexx46/esp32-backend-api/internal/mqtt"
	"github.com/rednexx46/esp32-backend-api/internal/ws"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	mqttLib "github.com/eclipse/paho.mqtt.golang"
)

// @title           ESP32 Backend API
// @version         1.0
// @description     REST API for handling sensor data, KPIs, and user authentication in an IoT ESP32 system.
// @termsOfService  https://github.com/rednexx46/esp32-backend-api

// @contact.name   José Xavier
// @contact.url    https://github.com/rednexx46
// @contact.email  josexavier46@outlook.pt

// @host      gateway.local:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] .env file not found, using environment variables")
	}

	// Initialize MongoDB
	db.InitDB()

	// Seed Admin User
	db.SeedAdminUser()

	// Start WebSocket hub
	go ws.StartHub()

	// Initialize MQTT client and subscribe to topic
	mqtt.InitMQTT(func(client mqttLib.Client, msg mqttLib.Message) {
		log.Printf("[MQTT] Received: %s", msg.Payload())
		ws.Broadcast(msg.Payload())
	})

	// Setup Gin router
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket live endpoint
	r.GET("/ws/live-data", ws.LiveDataWebSocket)

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/login", handlers.LoginHandler)
		public.POST("/logout", handlers.LogoutHandler)
	}

	// Protected routes
	protected := public.Group("/")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.GET("/profile", handlers.GetProfile)

		admin := protected.Group("/")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/data/:device_id", handlers.GetSensorDataByDevice)
			admin.GET("/data", handlers.GetAllSensorData)
			admin.GET("/devices", handlers.GetActiveDevices)
			admin.GET("/kpis", handlers.GetAllKPIs)
			admin.GET("/kpis/device/:device_id", handlers.GetKPIsByDevice)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[SERVER] Listening on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[SERVER] Failed to start: %v", err)
	}
}
