package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rednexx46/esp32-backend-api/internal/db"
	"github.com/rednexx46/esp32-backend-api/internal/handlers"
	"github.com/rednexx46/esp32-backend-api/internal/middleware"
	"github.com/rednexx46/esp32-backend-api/internal/models"
	"github.com/rednexx46/esp32-backend-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env")
	db.InitDB()
	seedTestUser()
	os.Exit(m.Run())
}

func seedTestUser() {
	hashed, _ := utils.HashPassword("testpass")
	user := models.User{
		Username:  "testuser",
		Password:  hashed,
		Role:      "admin",
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", handlers.LoginHandler)

	r.GET("/protected", middleware.JWTMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "access granted"})
	})

	r.GET("/admin-only", middleware.JWTMiddleware(), middleware.AdminOnly(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"admin": "true"})
	})

	return r
}

func TestLoginSuccess(t *testing.T) {
	r := setupRouter()

	body := models.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}

func TestLoginFailure(t *testing.T) {
	r := setupRouter()

	body := models.LoginRequest{
		Username: "testuser",
		Password: "wrongpass",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid credentials")
}

func TestProtectedEndpointWithToken(t *testing.T) {
	r := setupRouter()

	loginReq := models.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	jsonValue, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var result map[string]string
	json.Unmarshal(resp.Body.Bytes(), &result)
	token := result["token"]

	protectedReq, _ := http.NewRequest("GET", "/protected", nil)
	protectedReq.Header.Set("Authorization", "Bearer "+token)
	protectedResp := httptest.NewRecorder()
	r.ServeHTTP(protectedResp, protectedReq)

	assert.Equal(t, http.StatusOK, protectedResp.Code)
	assert.Contains(t, protectedResp.Body.String(), "access granted")
}

func TestProtectedEndpointWithoutToken(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/protected", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAdminOnlyEndpoint(t *testing.T) {
	r := setupRouter()

	// Login to get token
	loginReq := models.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	jsonValue, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var result map[string]string
	json.Unmarshal(resp.Body.Bytes(), &result)
	token := result["token"]

	// Access admin-only route
	adminReq, _ := http.NewRequest("GET", "/admin-only", nil)
	adminReq.Header.Set("Authorization", "Bearer "+token)
	adminResp := httptest.NewRecorder()
	r.ServeHTTP(adminResp, adminReq)

	assert.Equal(t, http.StatusOK, adminResp.Code)
	assert.Contains(t, adminResp.Body.String(), "admin")
}
