package repository

import (
	"context"
	"log"
	"os"

	"github.com/Madslick/chit-chat-go/internal/chat/datastructs"
	"github.com/Madslick/chit-chat-go/internal/chat/repository/mongostructs"
	"github.com/Madslick/chit-chat-go/internal/shared/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	conn db.DbConnection

	client        *mongo.Client
	conversations *mongo.Collection
}

func (mr *mongoRepository) init() {
	var err error
	mr.client, err = mr.conn.MongoClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	mr.conversations = mr.client.Database(os.Getenv("CHAT_DB")).Collection("conversations")
}

func (mr *mongoRepository) CreateConversation(members []*datastructs.Client) (string, error) {
	clients := []mongostructs.Client{}

	for _, mem := range members {
		clients = append(clients, mongostructs.Client{
			Id:   mem.ClientId,
			Name: mem.Name,
		})
	}
	conversation := mongostructs.Conversation{
		Members: clients,
	}

	result, err := mr.conversations.InsertOne(context.TODO(), conversation)
	if err != nil {
		log.Fatal("Error inserting new Conversation", err)
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).String(), nil
}

func (mr *mongoRepository) GetConversationByMemberNames(memberNames []string) (datastructs.Conversation, error) {

	memberFilter := bson.A{}
	for _, memberName := range memberNames {
		memberFilter = append(memberFilter, bson.M{"members": bson.M{"$elemMatch": bson.M{"name": memberName}}})
	}
	memberFilter = append(memberFilter, bson.M{"members": bson.M{"$size": len(memberNames)}})

	filter := bson.D{
		{"$and", memberFilter},
	}
	conversations := []mongostructs.Conversation{}
	cursor, err := mr.conversations.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &conversations); err != nil {
		panic(err)
	}

	log.Printf("Conversation filter returned %d results\n", len(conversations))

	conversation := datastructs.Conversation{}
	if len(conversations) == 1 {
		conversation.Id = conversations[0].Id.String()

		// Build Members
		conversation.Members = []*datastructs.Client{}
		for _, client := range conversations[0].Members {
			conversation.Members = append(conversation.Members,
				&datastructs.Client{
					ClientId: client.Id,
					Name:     client.Name,
				},
			)
		}
	}
	return conversation, nil
}

func (mr *mongoRepository) CreateMessage(message datastructs.Message) (string, error) {
	return "", nil
}
