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
var UserCollection = db.UserData(db.Client, "users")
var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity, adjust as needed
	},
}

func EnterApp(c *gin.Context) {
	userId := c.Query("user_id")
	fmt.Println("User with ID:", userId, " is entering the app")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("reached here 0")
	err := UserCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		fmt.Println("Error finding user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	fmt.Println("reached here 1")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	OnlineUsers[userId] = conn

	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	cursor, err := ConversationCollection.Find(ctx, bson.M{
		"participants.id": userId,
	})
	if err != nil {
		fmt.Println("Error fetching conversations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}
	defer cursor.Close(ctx)
	var conversations []models.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		fmt.Println("Error decoding conversations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode conversations"})
		return
	}
	type msg struct {
		UserId string `json:"user_id"`
		Online bool   `json:"online"`
	}
	for _, conversation := range conversations {
		for _, participants := range conversation.Participants {
			if participants.Id != userId {
				exists := OnlineUsers[participants.Id]
				if exists != nil {
					OtherUserConn := OnlineUsers[participants.Id]
					OtherUserConn.WriteJSON(&msg{
						UserId: userId,
						Online: true,
					})
					conn.WriteJSON(&msg{
						UserId: participants.Id,
						Online: true,
					})
				}
			}
		}
	}

	var exitMsg msg
	conn.ReadJSON(&exitMsg)
	var conversations1 []models.Conversation
	if err := cursor.All(ctx, &conversations1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode conversations"})
		return
	}
	for _, conversation := range conversations {
		for _, participants := range conversation.Participants {
			if participants.Id != userId {
				exists := OnlineUsers[participants.Id]
				if exists != nil {
					conn := OnlineUsers[participants.Id]
					conn.WriteJSON(&msg{
						UserId: userId,
						Online: false,
					})
				}
			}
		}
	}
	delete(OnlineUsers, userId)

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
