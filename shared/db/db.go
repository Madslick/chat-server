package db

type DbConnection interface {
	Connect() error
}

func SetupDb() DbConnection {
	dbConnection := &mongoClient{}
	dbConnection.Connect()
	return dbConnection
}
