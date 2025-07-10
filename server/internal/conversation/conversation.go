package conversation

import (
	"chat-server/db"
	user "chat-server/internal/users"
	"chat-server/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ConversationCollection = db.ConversationData(db.Client, "conversations")

func AddUserToConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		user_id, exits := c.Get("user_id")
		if !exits {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Please log in"})
			return
		}
		second_user_id := c.Param("user_id")
		if user_id == second_user_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot add yourself to a conversation"})
			return
		}
		var secondUser models.User
		err := user.UserCollection.FindOne(ctx, primitive.M{"user_id": second_user_id}).Decode(&secondUser)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": "Please check the user ID"})
			return
		}
		var conversation models.Conversation
		conversation.Id = primitive.NewObjectID()
		conversation.RoomId = conversation.Id.Hex()
		conversation.Participants = []string{user_id.(string), second_user_id}
		conversation.LastMessage = nil
		conversation.CreatedAt = time.Now()
		conversation.UpdatedAt = time.Now()
		_, err = ConversationCollection.InsertOne(ctx, conversation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation", "message": "Please try again later"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})

	}
}
func GetConversationByUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		user_id, exits := c.Get("user_id")
		if !exits {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Please log in"})
			return
		}
		var conversations []models.Conversation
		cursor, err := ConversationCollection.Find(ctx, primitive.M{"participants": user_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations", "message": err.Error()})
			return
		}
		if err = cursor.All(ctx, &conversations); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode conversations", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Conversations fetched successfully",
			"data":    conversations,
		})
	}
}
