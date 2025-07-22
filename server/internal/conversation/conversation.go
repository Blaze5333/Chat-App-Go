package conversation

import (
	"chat-server/db"
	user "chat-server/internal/users"
	"chat-server/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ConversationCollection = db.ConversationData(db.Client, "conversations")
var MessageCollection = db.MessageData(db.Client, "messages")

func AddUserToConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Participants: xndxndcindcek")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		user_id, exits := c.Get("user_id")
		if !exits {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Please log in"})
			return
		}
		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Please log in"})
			return
		}
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Please log in"})
			return
		}
		second_user_id := c.Param("user_id")
		if user_id == second_user_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot add yourself to a conversation"})
			return
		}

		var currentUser models.User
		err := user.UserCollection.FindOne(ctx, primitive.M{"user_id": user_id}).Decode(&currentUser)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Current user not found"})
			return
		}

		var secondUser models.User
		err = user.UserCollection.FindOne(ctx, primitive.M{"user_id": second_user_id}).Decode(&secondUser)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "message": "Please check the user ID"})
			return
		}
		if !secondUser.Verified {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User is not verified"})
			return
		}
		var conversation models.Conversation
		conversation.Id = primitive.NewObjectID()
		conversation.RoomId = conversation.Id.Hex()
		conversation.Participants = []models.Participant{
			{
				Id:       user_id.(string),
				Username: username.(string),
				Email:    email.(string),
				Image:    currentUser.Image,
			},
			{
				Id:       second_user_id,
				Username: secondUser.Username,
				Email:    secondUser.Email,
				Image:    secondUser.Image,
			},
		}

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
		log.Println("Fetching conversations for user")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		user_id, exits := c.Get("user_id")
		if !exits {
			log.Println("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Please log in"})
			return
		}
		log.Printf("User ID from context: %v", user_id)
		var conversations []models.Conversation
		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: bson.D{{Key: "participants.id", Value: user_id}}}},
			{{Key: "$unwind", Value: "$participants"}},
			{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "participants.id"},
				{Key: "foreignField", Value: "user_id"},
				{Key: "as", Value: "userInfo"},
			}}},
			{{Key: "$unwind", Value: "$userInfo"}},
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "participants", Value: bson.D{{Key: "$push", Value: bson.D{
					{Key: "id", Value: "$participants.id"},
					{Key: "username", Value: "$participants.username"},
					{Key: "image", Value: "$userInfo.image"},
				}}}},
			}}},
		}
		cursor, err := ConversationCollection.Aggregate(ctx, pipeline)
		if err != nil {
			log.Println("Error aggregating", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error", "error": err.Error()})
		}
		if err = cursor.All(ctx, &conversations); err != nil {
			log.Printf("Error decoding conversations: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode conversations", "message": err.Error()})
			return
		}
		log.Printf("Found %d conversations", len(conversations))
		c.JSON(http.StatusOK, gin.H{
			"message": "Conversations fetched successfully",
			"data":    conversations,
		})
	}
}

func GetRoomMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		roomId := c.Param("room_id")
		var messages []models.Message
		cursor, err := MessageCollection.Find(ctx, bson.M{"room_id": roomId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages", "message": err.Error()})
			return
		}
		defer cursor.Close(ctx)
		err = cursor.All(ctx, &messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode messages", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Messages fetched successfully",
			"data":    messages,
		})
	}
}
