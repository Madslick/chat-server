package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Madslick/chit-chat-go/internal/chat/datastruct"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	client        *mongo.Client
	conversations *mongo.Collection
}

func (mc *mongoClient) Connect() error {
	var err error
	mc.client, err = mongo.NewClient(
		options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s:%s@%s:%s",
				os.Getenv("MONGO_USER"),
				os.Getenv("MONGO_PASS"),
				os.Getenv("MONGO_HOST"),
				os.Getenv("MONGO_PORT"),
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mc.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := mc.client.Ping(ctx, nil); err != nil {
		log.Fatal("Unable to ping mongodb with newly created connection", err)
		return err
	}
	log.Println("MongoDB Connected Successfully")

	mc.conversations = mc.client.Database(os.Getenv("CHAT_DB")).Collection("conversations")

	return nil
}

func (mc *mongoClient) CreateConversation(members []*datastruct.Client) {

}

func (mc *mongoClient) CreateMessage(message *datastruct.Message) {

}
