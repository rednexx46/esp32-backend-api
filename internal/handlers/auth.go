package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rednexx46/esp32-backend-api/internal/db"
	"github.com/rednexx46/esp32-backend-api/internal/models"
	"github.com/rednexx46/esp32-backend-api/internal/utils"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GetProfile godoc
// @Summary      Get authenticated user's profile
// @Description  Retrieves the profile information of the currently authenticated user.
// @Tags         auth
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "User profile"
// @Failure      401  {object}  map[string]string       "Unauthorized"
// @Router       /auth/profile [get]
func GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// LoginHandler godoc
// @Summary      Login user
// @Description  Authenticates a user and returns a JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login  body      models.LoginRequest  true  "Login credentials"
// @Success      200    {object}  models.LoginResponse "JWT token"
// @Failure      400    {object}  map[string]string    "Invalid request body"
// @Failure      401    {object}  map[string]string    "Invalid username or password"
// @Failure      500    {object}  map[string]string    "Failed to generate token"
// @Router       /auth/login [post]
func LoginHandler(c *gin.Context) {
	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx := context.Background()
	user, err := db.FindUserByUsername(ctx, login.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !utils.ComparePassword(login.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	TOKEN_TTL_STR := os.Getenv("TOKEN_TTL_MINUTES")
	TOKEN_TTL := 30 // Default to 30 minutes if not set
	if TOKEN_TTL_STR != "" {
		if ttl, err := strconv.Atoi(TOKEN_TTL_STR); err == nil {
			TOKEN_TTL = ttl
		}
	}
	expiration := time.Now().Add(time.Minute * time.Duration(TOKEN_TTL))
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := models.LoginResponse{
		Token: tokenString,
	}
	c.JSON(http.StatusOK, response)
}

// LogoutHandler handles user logout by revoking the authentication token.
//
// @Summary      Logout user
// @Description  Logs out the currently authenticated user by revoking their token.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string  "Logged out successfully"
// @Failure      401  {object}  map[string]string  "Unauthorized"
// @Router       /auth/logout [post]
func LogoutHandler(c *gin.Context) {
	// TODO: implement token revocation
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
