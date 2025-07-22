package ws

import (
	"chat-server/db"
	"chat-server/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var ConversationCollection = db.ConversationData(db.Client, "conversations")
var MessageCollection = db.MessageData(db.Client, "messages")
var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity, adjust as needed
	},
}

func getConversationByRoomId(roomId string) (*models.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var conversation models.Conversation
	err := ConversationCollection.FindOne(ctx, bson.M{"room_id": roomId}).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (h *Hub) HandleJoinRoom(c *gin.Context) {
	roomId := c.Param("room_id")
	userId := c.Query("user_id")
	userName := c.Query("username")
	_, err := getConversationByRoomId(roomId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	fmt.Print("User joined room:", userId, "Room ID:", roomId, "Username:", userName, "\n")
	client := &Client{
		ID:       userId,
		Conn:     conn,
		Message:  make(chan *models.Message),
		RoomId:   roomId,
		Username: userName,
	}
	fmt.Println("reached 1")
	_, exist := Rooms[roomId]
	if !exist {
		Rooms[roomId] = &Room{
			ID:      roomId,
			Clients: make(map[string]*Client),
		}
		fmt.Println("Client coming", client.ID, "in room", roomId)
		Rooms[roomId].Clients[client.ID] = client
	} else {
		Rooms[roomId].Clients[client.ID] = client
	}
	go client.WriteMessage()
	client.ReadMessage(h)
}
