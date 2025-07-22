package ws

import (
	"chat-server/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID       string `json:"id"`
	Conn     *websocket.Conn
	Message  chan *models.Message
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
}

var Rooms = make(map[string]*Room)

type Message models.Message
type Room struct {
	ID      string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

func (cl *Client) WriteMessage() {
	defer func() {
		cl.Conn.Close()
	}()
	for {
		msg, ok := <-cl.Message
		if !ok {
			return
		}
		fmt.Println("Sending message to client:", cl.ID, "in room:", cl.RoomId, "Content:", msg.Content)
		if err := cl.Conn.WriteJSON(msg); err != nil {
			return
		}
	}
}
func (cl *Client) ReadMessage(hub *Hub) {
	defer func() {
		fmt.Println("------Closing connection for client-------:", cl.ID)
		hub.Unregister <- cl
		cl.Conn.Close()
	}()

	for {
		_, msg, err := cl.Conn.ReadMessage()

		if err != nil {
			return
		}

		// Create a new context for each message insertion
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		log.Println("Received message from client:", cl.ID, "in room:", cl.RoomId, "Content:", string(msg))
		userMessage := &models.Message{
			Id:        primitive.NewObjectID(),
			RoomId:    cl.RoomId,
			Content:   string(msg),
			Username:  cl.Username,
			UserId:    cl.ID,
			CreatedAt: time.Now(),
		}

		_, err = MessageCollection.InsertOne(ctx, userMessage)
		cancel() // Cancel the context after the operation completes

		if err != nil {
			log.Println("Error inserting message:", err)
			return
		}

		hub.Broadcast <- userMessage

	}
}
