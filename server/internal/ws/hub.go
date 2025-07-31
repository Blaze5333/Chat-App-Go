package ws

import (
	"chat-server/models"
	"fmt"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *models.Message
}

var OnlineUsers = make(map[string]*websocket.Conn)

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *models.Message),
	}
}

type Notification struct {
	UserId   string `json:"user_id"`
	Content  string `json:"content"`
	Username string `json:"username"`
	Type     string `json:"type"` // "message" or "notification"
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			room, exists := h.Rooms[client.RoomId]
			fmt.Println("Client joined room", client.ID)
			if !exists {
				continue
			}
			_, exists = room.Clients[client.ID]
			if exists {
				continue
			}
			room.Clients[client.ID] = client
		case client := <-h.Unregister:
			room, exists := h.Rooms[client.RoomId]
			if !exists {
				continue
			}
			_, exists = room.Clients[client.ID]
			if !exists {
				continue
			}
			delete(h.Rooms[client.RoomId].Clients, client.ID)
			close(client.Message)
		case message := <-h.Broadcast:
			for _, room := range Rooms {
				if room.ID == message.RoomId {
					for _, client := range room.Clients {
						client.Message <- message
						if client.ID != message.UserId {
							exists := OnlineUsers[client.ID]
							if exists != nil {
								conn := OnlineUsers[client.ID]
								conn.WriteJSON(&Notification{
									UserId:   message.UserId,
									Content:  message.Content,
									Username: message.Username,
									Type:     "notification",
								})
							}
						}
					}
				}

			}

		}
	}
}
