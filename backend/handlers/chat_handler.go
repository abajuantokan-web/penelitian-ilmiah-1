package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

type Client struct {
	conn	*websocket.Conn
	mu	sync.Mutex
}

func (c *Client) safeWrite(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

type Hub struct {
	mu		sync.RWMutex
	connections	map[int32]*Client
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[int32]*Client),
	}
}
func (h *Hub) Register(userID int32, conn *websocket.Conn) *Client {
	h.mu.Lock()
	defer h.mu.Unlock()

	if existing, ok := h.connections[userID]; ok {
		existing.conn.Close()
	}
	client := &Client{conn: conn}
	h.connections[userID] = client
	log.Printf("🔗 User %d connected to chat", userID)
	return client
}

func (h *Hub) Unregister(userID int32, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if c, ok := h.connections[userID]; ok && c == client {
		c.conn.Close()
		delete(h.connections, userID)
		log.Printf("🔌 User %d disconnected from chat", userID)
	}
}

func (h *Hub) SendToUser(userID int32, message interface{}) bool {
	h.mu.RLock()

	client, ok := h.connections[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("⚠️ Failed to marshal message for user %d: %v", userID, err)
		return false
	}

	if err := client.safeWrite(websocket.TextMessage, data); err != nil {
		log.Printf("⚠️ Failed to send message to user %d: %v", userID, err)
		return false
	}

	return true
}

var ChatHub = NewHub()

var upgrader = websocket.Upgrader{
	ReadBufferSize:		1024,
	WriteBufferSize:	1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {

	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "token is required"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "my_super_secret_key_for_openpeo_platform" // fallback lokal
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token claims"})
		return
	}

	senderIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid user_id in token"})
		return
	}
	senderID := int32(senderIDFloat)

	receiverIDStr := c.Query("receiver_id")
	var receiverID int32
	if receiverIDStr != "" {
		id, err := strconv.ParseInt(receiverIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":	false,
				"message":	"Invalid receiver_id",
			})
			return
		}
		receiverID = int32(id)
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade failed: %v", err)
		return
	}

	client := ChatHub.Register(senderID, conn)

	defer func() {
		ChatHub.Unregister(senderID, client)
	}()

	fmt.Printf("💬 Chat channel opened: User %d (receiver: %d)\n", senderID, receiverID)

	const pongWait = 60 * time.Second
	const pingPeriod = 54 * time.Second

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		for {
			<-ticker.C
			if err := client.safeWrite(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

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

		payload.SenderID = senderID
		if payload.ReceiverID == 0 {
			payload.ReceiverID = receiverID
		}

		var senderUser models.User
		if config.DB.Where("id = ?", payload.SenderID).First(&senderUser).Error == nil && senderUser.Role == "customer" {
			var targetUser models.User
			if err := config.DB.Where("id = ?", payload.ReceiverID).First(&targetUser).Error; err != nil || (targetUser.Role != "seller" && targetUser.Role != "admin" && targetUser.Role != "vendor") {
				var sp models.SellerProfile
				if err := config.DB.Where("id = ?", payload.ReceiverID).First(&sp).Error; err == nil {
					payload.ReceiverID = sp.UserID
				}
			}
		}

		message := models.Message{
			SenderID:	payload.SenderID,
			ReceiverID:	payload.ReceiverID,
			Content:	payload.Content,
			IsRead:		false,
		}

		if err := config.DB.Create(&message).Error; err != nil {
			log.Printf("❌ Failed to save message to database: %v", err)
			continue
		}

		response := gin.H{
			"id":		message.ID,
			"sender_id":	message.SenderID,
			"receiver_id":	message.ReceiverID,
			"content":	message.Content,
			"is_read":	message.IsRead,
			"created_at":	message.CreatedAt,
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
			"success":	false,
			"message":	"sender_id and receiver_id are required",
		})
		return
	}

	var senderUser models.User
	if config.DB.Where("id = ?", senderID).First(&senderUser).Error == nil && senderUser.Role == "customer" {
		if id, err := strconv.Atoi(receiverID); err == nil {
			var targetUser models.User
			if err := config.DB.Where("id = ?", id).First(&targetUser).Error; err != nil || (targetUser.Role != "seller" && targetUser.Role != "admin" && targetUser.Role != "vendor") {
				var sp models.SellerProfile
				if err := config.DB.Where("id = ?", id).First(&sp).Error; err == nil {
					receiverID = strconv.Itoa(int(sp.UserID))
				}
			}
		}
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
			"success":	false,
			"message":	"Failed to fetch chat history",
			"error":	result.Error.Error(),
		})
		return
	}

	config.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", receiverID, senderID, false).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"data":		messages,
		"total":	len(messages),
	})
}

func GetChatContacts(c *gin.Context) {
	userIDStr := c.Query("user_id")
	role := c.Query("role")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user_id is required"})
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	var contacts []map[string]interface{}

	if role == "seller" {
		query := `
		SELECT u.id, u.name, 
			(SELECT COUNT(*) FROM messages m WHERE m.sender_id = u.id AND m.receiver_id = ? AND m.is_read = 0) as unread_count,
			(SELECT content FROM messages m2 WHERE (m2.sender_id = u.id AND m2.receiver_id = ?) OR (m2.receiver_id = u.id AND m2.sender_id = ?) ORDER BY m2.created_at DESC LIMIT 1) as last_message,
			MAX(m.created_at) as last_activity
		FROM users u
		JOIN messages m ON (m.sender_id = u.id AND m.receiver_id = ?) OR (m.receiver_id = u.id AND m.sender_id = ?)
		WHERE u.role = 'customer'
		GROUP BY u.id, u.name
		ORDER BY last_activity DESC
		`
		config.DB.Raw(query, userID, userID, userID, userID, userID).Scan(&contacts)
	} else {
		query := `
		SELECT u.id, COALESCE(NULLIF(u.store_name, ''), u.name) as name,
			(SELECT COUNT(*) FROM messages m WHERE m.sender_id = u.id AND m.receiver_id = ? AND m.is_read = 0) as unread_count,
			(SELECT content FROM messages m2 WHERE (m2.sender_id = u.id AND m2.receiver_id = ?) OR (m2.receiver_id = u.id AND m2.sender_id = ?) ORDER BY m2.created_at DESC LIMIT 1) as last_message,
			MAX(m.created_at) as last_activity
		FROM users u
		JOIN messages m ON (m.sender_id = u.id AND m.receiver_id = ?) OR (m.receiver_id = u.id AND m.sender_id = ?)
		WHERE u.role IN ('seller', 'admin', 'vendor')
		GROUP BY u.id, u.name, u.store_name
		ORDER BY last_activity DESC
		`
		config.DB.Raw(query, userID, userID, userID, userID, userID).Scan(&contacts)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": contacts})
}
