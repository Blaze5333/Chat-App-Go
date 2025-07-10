package routes

import (
	"chat-server/internal/conversation"
	user "chat-server/internal/users"
	"chat-server/internal/ws"
	"chat-server/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/register", user.RegisterUser())
	incomingRoutes.POST("/login", user.LoginUser())
}
func ChatRoutes(incomingRoutes *gin.Engine, ws *ws.Hub) {
	incomingRoutes.POST("/create_room/:user_id", middleware.Authenticate(), conversation.AddUserToConversation())
	incomingRoutes.GET("/conversation", middleware.Authenticate(), conversation.GetConversationByUserId())
	incomingRoutes.GET("/join_room/:room_id", ws.HandleJoinRoom)
}
