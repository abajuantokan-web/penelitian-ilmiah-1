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







type Hub struct {
	mu          sync.RWMutex
	connections map[int32]*websocket.Conn
}


func NewHub() *Hub {
	return &Hub{
		connections: make(map[int32]*websocket.Conn),
	}
}


func (h *Hub) Register(userID int32, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	
	if existing, ok := h.connections[userID]; ok {
		existing.Close()
	}
	h.connections[userID] = conn
	log.Printf("🔗 User %d connected to chat", userID)
}


func (h *Hub) Unregister(userID int32) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.connections[userID]; ok {
		conn.Close()
		delete(h.connections, userID)
		log.Printf("🔌 User %d disconnected from chat", userID)
	}
}



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






var ChatHub = NewHub()



var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}



















func HandleWebSocket(c *gin.Context) {
	
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

	
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade failed: %v", err)
		return
	}

	
	ChatHub.Register(int32(senderID), conn)

	
	defer func() {
		ChatHub.Unregister(int32(senderID))
	}()

	fmt.Printf("💬 Chat channel opened: User %d (receiver: %d)\n", senderID, receiverID)

	
	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("⚠️ WebSocket error for user %d: %v", senderID, err)
			}
			break
		}

		
		var payload models.ChatPayload
		if err := json.Unmarshal(rawMessage, &payload); err != nil {
			log.Printf("⚠️ Invalid message format from user %d: %v", senderID, err)
			continue
		}

		
		payload.SenderID = int32(senderID)
		if payload.ReceiverID == 0 {
			payload.ReceiverID = int32(receiverID)
		}

		
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

		
		response := gin.H{
			"id":          message.ID,
			"sender_id":   message.SenderID,
			"receiver_id": message.ReceiverID,
			"content":     message.Content,
			"is_read":     message.IsRead,
			"created_at":  message.CreatedAt,
		}

		
		ChatHub.SendToUser(payload.ReceiverID, response)

		
		ChatHub.SendToUser(payload.SenderID, response)
	}
}












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

	
	config.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", receiverID, senderID, false).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    messages,
		"total":   len(messages),
	})
}
