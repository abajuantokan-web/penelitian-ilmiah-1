package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

// -----------------------------------------------------------------------
// WebSocket Hub — manages active connections for real-time chat routing
// -----------------------------------------------------------------------

// Hub maintains a registry of active WebSocket connections indexed by user ID.
// It provides thread-safe connection registration, removal, and message routing.
type Hub struct {
	mu          sync.RWMutex
	connections map[int32]*websocket.Conn
}

// NewHub creates and returns a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		connections: make(map[int32]*websocket.Conn),
	}
}

// Register adds a user's WebSocket connection to the hub.
func (h *Hub) Register(userID int32, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Close any existing connection for this user (prevents stale sockets)
	if existing, ok := h.connections[userID]; ok {
		existing.Close()
	}
	h.connections[userID] = conn
	log.Printf("🔗 User %d connected to chat", userID)
}

// Unregister removes a user's connection from the hub.
func (h *Hub) Unregister(userID int32) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.connections[userID]; ok {
		conn.Close()
		delete(h.connections, userID)
		log.Printf("🔌 User %d disconnected from chat", userID)
	}
}

// SendToUser delivers a JSON message to a specific user if they are online.
// Returns true if the message was delivered, false if the user is offline.
func (h *Hub) SendToUser(userID int32, message interface{}) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, ok := h.connections[userID]
	if !ok {
		return false
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("⚠️ Failed to marshal message for user %d: %v", userID, err)
		return false
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("⚠️ Failed to send message to user %d: %v", userID, err)
		return false
	}

	return true
}

// -----------------------------------------------------------------------
// Global Hub instance and WebSocket upgrader
// -----------------------------------------------------------------------

// ChatHub is the singleton hub instance used across the application.
var ChatHub = NewHub()

// upgrader configures the HTTP-to-WebSocket upgrade.
// CheckOrigin is permissive for development; restrict in production.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// -----------------------------------------------------------------------
// WebSocket Handler
// -----------------------------------------------------------------------

// HandleWebSocket handles GET /ws/chat
// It upgrades the HTTP connection to a full-duplex WebSocket and manages
// real-time bidirectional chat between users.
//
// Query Parameters:
//   - sender_id:   The ID of the user initiating the connection (required)
//   - receiver_id: The ID of the intended chat partner (required)
//
// Message Flow:
//  1. Client connects with sender_id and receiver_id
//  2. Connection is registered in the Hub under sender_id
//  3. Incoming messages are persisted to MySQL immediately
//  4. Messages are routed to receiver's connection if online
//  5. On disconnect, the connection is cleaned up from the Hub
func HandleWebSocket(c *gin.Context) {
	// Parse sender and receiver IDs from query parameters
	senderIDStr := c.Query("sender_id")
	receiverIDStr := c.Query("receiver_id")

	if senderIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "sender_id is a required query parameter",
		})
		return
	}

	senderID, err := strconv.ParseInt(senderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid sender_id",
		})
		return
	}

	var receiverID int64
	if receiverIDStr != "" {
		receiverID, err = strconv.ParseInt(receiverIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid receiver_id",
			})
			return
		}
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade failed: %v", err)
		return
	}

	// Register this connection in the hub
	ChatHub.Register(int32(senderID), conn)

	// Ensure cleanup on disconnect
	defer func() {
		ChatHub.Unregister(int32(senderID))
	}()

	fmt.Printf("💬 Chat channel opened: User %d (receiver: %d)\n", senderID, receiverID)

	// Read loop — continuously listen for incoming messages
	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("⚠️ WebSocket error for user %d: %v", senderID, err)
			}
			break
		}

		// Parse the incoming chat payload
		var payload models.ChatPayload
		if err := json.Unmarshal(rawMessage, &payload); err != nil {
			log.Printf("⚠️ Invalid message format from user %d: %v", senderID, err)
			continue
		}

		// Ensure sender_id consistency
		payload.SenderID = int32(senderID)
		if payload.ReceiverID == 0 {
			payload.ReceiverID = int32(receiverID)
		}

		// Persist message to MySQL immediately
		message := models.Message{
			SenderID:   payload.SenderID,
			ReceiverID: payload.ReceiverID,
			Content:    payload.Content,
			IsRead:     false,
		}

		if err := config.DB.Create(&message).Error; err != nil {
			log.Printf("❌ Failed to save message to database: %v", err)
			continue
		}

		// Prepare the response payload with database-generated fields
		response := gin.H{
			"id":          message.ID,
			"sender_id":   message.SenderID,
			"receiver_id": message.ReceiverID,
			"content":     message.Content,
			"is_read":     message.IsRead,
			"created_at":  message.CreatedAt,
		}

		// Route message to receiver if they are online
		ChatHub.SendToUser(payload.ReceiverID, response)

		// Echo confirmation back to sender
		ChatHub.SendToUser(payload.SenderID, response)
	}
}

// -----------------------------------------------------------------------
// Chat History REST Endpoint
// -----------------------------------------------------------------------

// GetChatHistory handles GET /api/messages
// Retrieves the conversation history between two users.
//
// Query Parameters:
//   - sender_id:   First participant's ID (required)
//   - receiver_id: Second participant's ID (required)
//   - limit:       Max messages to return (default: 50)
func GetChatHistory(c *gin.Context) {
	senderID := c.Query("sender_id")
	receiverID := c.Query("receiver_id")

	if senderID == "" || receiverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "sender_id and receiver_id are required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 || limit > 200 {
		limit = 50
	}

	var messages []models.Message

	// Fetch messages in both directions between the two users
	result := config.DB.
		Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID,
		).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch chat history",
			"error":   result.Error.Error(),
		})
		return
	}

	// Mark unread messages as read for the current user (receiver of these messages)
	config.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", receiverID, senderID, false).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    messages,
		"total":   len(messages),
	})
}
