package routes

import (
	"chat-server/internal/conversation"
	user "chat-server/internal/users"
	"chat-server/internal/ws"
	"chat-server/middleware"

	"github.com/gin-gonic/gin"
)

func SocialRRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/auth/:provider", user.SocialLogin())
	incomingRoutes.GET("/auth/:provider/redirect", user.SocialLoginCallback())
}
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/register", user.RegisterUser())
	incomingRoutes.POST("/login", user.LoginUser())
	incomingRoutes.GET("/users/search" /*middleware.Authenticate(),*/, user.GetUserByEmail())
	incomingRoutes.POST("/verify_otp", user.VerifyOtp())
	incomingRoutes.POST("/upload_image", middleware.Authenticate(), user.UploadHandler)
}

func ChatRoutes(incomingRoutes *gin.Engine, ws *ws.Hub) {
	incomingRoutes.POST("/create_room/:user_id", middleware.Authenticate(), conversation.AddUserToConversation())
	incomingRoutes.GET("/conversation", middleware.Authenticate(), conversation.GetConversationByUserId())
	incomingRoutes.GET("/join_room/:room_id", ws.HandleJoinRoom)
	incomingRoutes.GET("/get_room_messages/:room_id", middleware.Authenticate(), conversation.GetRoomMessages())
}
