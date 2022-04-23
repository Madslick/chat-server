package repository

import (
	"context"
	"log"
	"os"

	"github.com/Madslick/chit-chat-go/internal/auth/datastructs"
	"github.com/Madslick/chit-chat-go/internal/auth/repository/mongostructs"
	"github.com/Madslick/chit-chat-go/internal/shared/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	conn db.DbConnection

	client   *mongo.Client
	accounts *mongo.Collection
}

func (mr *mongoRepository) init() {
	var err error
	mr.client, err = mr.conn.MongoClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	mr.accounts = mr.client.Database(os.Getenv("CHAT_DB")).Collection("accounts")
}

func (mr *mongoRepository) SignUp(email string, password string, first string, last string, phone string) (string, error) {
	result, err := mr.accounts.InsertOne(
		context.TODO(),
		mongostructs.Account{
			Email:    email,
			Password: password,
			First:    first,
			Last:     last,
			Phone:    phone,
		},
	)
	if err != nil {
		panic(err)
	}

	return result.InsertedID.(primitive.ObjectID).String(), nil
}

func (mr *mongoRepository) SignIn(email string, password string) (*datastructs.Account, error) {
	var account mongostructs.Account
	cursor := mr.accounts.FindOne(
		context.TODO(),
		bson.M{
			"email":    email,
			"password": password,
		},
	)

	if err := cursor.Decode(&account); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &datastructs.Account{
		Id:    account.Id,
		Email: account.Email,
		First: account.First,
		Last:  account.Last,
		Phone: account.Phone,
	}, nil
}

func (mr *mongoRepository) SearchAccounts(searchQuery string) ([]*datastructs.Account, error) {

	return []*datastructs.Account{}, nil
}
