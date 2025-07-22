package ws

import (
	"chat-server/models"
	"fmt"
)

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *models.Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *models.Message),
	}
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
					}
				}

			}

		}
	}
}
