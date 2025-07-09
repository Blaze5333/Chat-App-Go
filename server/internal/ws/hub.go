package ws

import "fmt"

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
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
			room, exists := h.Rooms[message.RoomId]
			if !exists {
				continue
			}
			for _, client := range room.Clients {
				fmt.Println("room clients", client)
				client.Message <- message
			}

		}
	}
}
