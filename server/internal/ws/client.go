package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string `json:"id"`
	Conn     *websocket.Conn
	Message  chan *Message
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
}
type Message struct {
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}
type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
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
		if err := cl.Conn.WriteJSON(msg); err != nil {
			return
		}
	}
}
func (cl *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- cl
		cl.Conn.Close()

	}()
	for {
		_, msg, err := cl.Conn.ReadMessage()
		fmt.Println("Received message:", string(msg))
		if err != nil {
			return
		}
		message := &Message{
			RoomId:   cl.RoomId,
			Username: cl.Username,
			Content:  string(msg),
		}
		hub.Broadcast <- message

	}
}
