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
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Build conversation mongostruct
	clients := []mongostructs.Client{}
	for _, mem := range members {
		clients = append(clients, mongostructs.Client{
			Id:   mem.ClientId,
			Name: mem.Name,
		})
	}
	conversation := mongostructs.Conversation{
		Members:  clients,
		Messages: make([]mongostructs.Message, 0),
	}

	// Insert into database and retrieve ID
	result, err := mr.conversations.InsertOne(context.TODO(), conversation)
	if err != nil {
		log.Fatal("Error inserting new Conversation", err)
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (mr *mongoRepository) GetConversationByMemberIds(memberIds []string) (datastructs.Conversation, error) {

	// Build Filter object
	memberFilter := bson.A{}
	for _, memberId := range memberIds {
		memberFilter = append(
			memberFilter,
			bson.M{"members": bson.M{"$elemMatch": bson.M{"id": memberId}}},
		)
	}
	memberFilter = append(memberFilter, bson.M{"members": bson.M{"$size": len(memberIds)}})
	filter := bson.D{
		{"$and", memberFilter},
	}

	opts := options.Find().SetProjection(bson.M{"messages": bson.M{"$slice": bson.A{0, 5}}})

	// Query database with filter
	conversations := []mongostructs.Conversation{}
	cursor, err := mr.conversations.Find(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}

	// Set result to conversations variable
	if err = cursor.All(context.TODO(), &conversations); err != nil {
		panic(err)
	}
	log.Printf("Conversation filter returned %d results\n", len(conversations))

	// Convert result to datastructs
	conversation := datastructs.Conversation{}
	if len(conversations) > 0 {
		conversation.Id = conversations[0].Id.Hex()

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

		// Build the history of messages
		conversation.Messages = []*datastructs.Message{}
		var msg mongostructs.Message
		for _, msg = range conversations[0].Messages {
			conversation.Messages = append(conversation.Messages, &datastructs.Message{
				From: datastructs.Client{
					ClientId: msg.From.Id,
					Name:     msg.From.Name,
				},
				Content: msg.Content,
			})
		}

	}

	return conversation, nil
}

func (mr *mongoRepository) CreateMessage(convId string, message datastructs.Message) (bool, error) {

	newMessage := mongostructs.Message{
		From: mongostructs.Client{
			Id:   message.From.ClientId,
			Name: message.From.Name,
		},
		Content: message.Content,
	}
	conversationId, _ := primitive.ObjectIDFromHex(convId)
	log.Printf("Create Message requested for %v\n", conversationId)
	result, err := mr.conversations.UpdateByID(
		context.TODO(),
		conversationId,
		bson.M{
			"$push": bson.M{
				"messages": newMessage,
			},
		},
	)
	if err != nil {
		log.Fatalf("Error occurred saving new message to conversation %s\n", conversationId)
		return false, err
	}

	updated := result.ModifiedCount == 1
	if updated {
		log.Printf("Message %s saved to db\n", newMessage.Content)
	}
	return updated, nil
}
