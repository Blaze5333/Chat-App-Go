package user

import (
	"chat-server/db"
	"chat-server/models"
	"chat-server/tokens"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection = db.UserData(db.Client, "users")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic("Error hashing password:", err)
		return "", err
	}
	return string(bytes), nil
}
func VerifyPassword(userpassword, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.UserRegisterReq
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err, "message": "Invalid request"})
			return
		}
		count, err := UserCollection.CountDocuments(ctx, primitive.M{"email": user.Email})
		if err != nil {
			c.JSON(500, gin.H{"error": err, "message": "Internal server error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists", "message": "Please use a different email"})
			return
		}
		var userData models.User
		userData.Username = user.Username
		userData.ID = primitive.NewObjectID()
		userData.Email = user.Email
		userData.UserId = userData.ID.Hex()

		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "Error hashing password"})
			return
		}
		userData.Password = hashedPassword
		_, err = UserCollection.InsertOne(ctx, userData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "Error inserting user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": gin.H{
			"id":       userData.UserId,
			"username": userData.Username,
			"email":    userData.Email,
		}})

	}
}
func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.UserLoginReq
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err, "message": "Invalid request"})
			return
		}
		var foundUser models.User
		err := UserCollection.FindOne(ctx, primitive.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "message": "Please check your credentials"})
			return
		}
		isValid, msg := VerifyPassword(user.Password, foundUser.Password)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg, "message": "Invalid password"})
			return
		}
		// Generate JWT token here (not implemented in this snippet)
		token, err := tokens.GenerateToken(foundUser.Email, foundUser.UserId, foundUser.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": gin.H{
			"id":       foundUser.UserId,
			"username": foundUser.Username,
			"email":    foundUser.Email,
			"token":    token,
		}})
	}
}
func GetUserByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		email := c.Param("email")
		var user models.User
		err := UserCollection.FindOne(ctx, primitive.M{"email": email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": "Please check the email"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"email": user.Email, "username": user.Username, "id": user.UserId})
	}
}

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
		err := UserCollection.FindOne(ctx, primitive.M{"user_id": second_user_id}).Decode(&secondUser)
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

	}
}
