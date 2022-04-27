package mongostructs

type Message struct {
	From    Client `bson:"client,omitempty"`
	Content string `bson:"content,omitempty"`
}
