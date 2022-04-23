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
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (mr *mongoRepository) SearchAccounts(searchQuery string, page int64, size int64) ([]*datastructs.Account, error) {
	opts := options.Find().SetSkip(page * size)
	cursor, err := mr.accounts.Find(
		context.TODO(),
		bson.M{
			"first": bson.M{
				"$regex": searchQuery,
			},
		},
		opts,
	)

	if err != nil {
		log.Fatalf("Failed to search accounts in mongodb %s\n", err)
		return nil, err
	}

	var results []mongostructs.Account
	err = cursor.All(context.TODO(), &results)

	if err != nil {
		log.Fatalf("Error occurred while iterating over the cursor of the search accounts result of mongodb %s\n", err)
		return nil, err
	}
	log.Printf("SearchAccounts found %d records\n", len(results))

	var accounts []*datastructs.Account
	for _, acc := range results {
		accounts = append(accounts, &datastructs.Account{
			Id:    acc.Id,
			Email: acc.Email,
			Phone: acc.Phone,
			First: acc.First,
			Last:  acc.Last,
		})
	}

	return accounts, nil
}
