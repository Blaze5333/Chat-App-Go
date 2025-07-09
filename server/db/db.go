package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return nil
	}
	var DATABASE_URI = os.Getenv("DATABASE_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(DATABASE_URI))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return nil
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return nil
	}
	log.Println("Connected to MongoDB successfully")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Chat_App").Collection(collectionName)
	return collection
}
func ConversationData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Chat_App").Collection(collectionName)
	return collection
}
func MessageData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Chat_App").Collection(collectionName)
	return collection
}
