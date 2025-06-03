package ws

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

var (
	clients   = make(map[*Client]bool)
	mutex     sync.Mutex
	broadcast = make(chan []byte)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// StartHub initializes the broadcast loop
func StartHub() {
	go func() {
		for msg := range broadcast {
			mutex.Lock()
			for client := range clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					client.Conn.Close()
					delete(clients, client)
				}
			}
			mutex.Unlock()
		}
	}()
}

// AddClient adds a new client connection to the list
func AddClient(conn *websocket.Conn) {
	client := &Client{Conn: conn}
	mutex.Lock()
	clients[client] = true
	mutex.Unlock()
}

// Broadcast sends a message to all connected clients
func Broadcast(msg []byte) {
	broadcast <- msg
}

// LiveDataWebSocket godoc
// @Summary      WebSocket real-time data stream
// @Description  Opens a WebSocket connection to receive real-time sensor data pushed from the backend. Requires a valid JWT token in the Authorization header.
// @Tags         websocket
// @Produce      json
// @Security     BearerAuth
// @Success      101  {string}  string  "Switching Protocols"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /ws/live-data [get]
func LiveDataWebSocket(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed Authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	AddClient(conn)
}
