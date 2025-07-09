package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Hub) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if _, exists := h.Rooms[req.ID]; exists {
		c.JSON(400, gin.H{"error": "Room already exists"})
		return
	}
	room := &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	h.Rooms[req.ID] = room
	c.JSON(201, room)
}

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity, adjust as needed
	},
}

func (h *Hub) HandleJoinRoom(c *gin.Context) {
	fmt.Println("here i am")

	roomId := c.Param("room_id")
	userId := c.Query("user_id")
	userName := c.Query("username")
	_, exists := h.Rooms[roomId]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	client := &Client{
		ID:       userId,
		Conn:     conn,
		Message:  make(chan *Message),
		RoomId:   roomId,
		Username: userName,
	}
	m := &Message{
		Content:  "New user joined the room",
		RoomId:   roomId,
		Username: userName,
	}
	h.Register <- client
	h.Broadcast <- m
	go client.WriteMessage()
	client.ReadMessage(h)
}
