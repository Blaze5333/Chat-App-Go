package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	UserId   string             `json:"user_id" bson:"user_id"`
}
type UserRegisterReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserLoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Conversation struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	Participants []string           `json:"participants" bson:"participants"`
	LastMessage  *Message           `json:"last_message" bson:"last_message"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	RoomId       string             `json:"room_id" bson:"room_id"`
}
type Message struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	RoomId    string             `json:"room_id" bson:"room_id"`
	Username  string             `json:"username" bson:"username"`
	Content   string             `json:"content" bson:"content"`
	UserId    string             `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
