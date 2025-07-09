package routes

import (
	"chat-server/internal/ws"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(incomingRoutes *gin.Engine, ws *ws.Hub) {
	incomingRoutes.POST("/create_room", ws.CreateRoom)
	incomingRoutes.GET("/join_room/:room_id", ws.HandleJoinRoom)
}
