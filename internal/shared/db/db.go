package db

import "go.mongodb.org/mongo-driver/mongo"

type DbConnection interface {
	MongoClient() (*mongo.Client, error)
	Connect() error
}

func SetupDb() DbConnection {
	dbConnection := &mongoClient{}
	dbConnection.Connect()
	return dbConnection
}
