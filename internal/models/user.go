package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"` // Do not expose password hash in JSON
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// LoginRequest represents the payload to request a login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the payload sent back after successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthClaims represents JWT claims.
type AuthClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
}

// TokenPayload represents stored token metadata.
type TokenPayload struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
