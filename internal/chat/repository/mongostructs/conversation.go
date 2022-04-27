package mongostructs

import "go.mongodb.org/mongo-driver/bson/primitive"

type Conversation struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Members  []Client           `bson:"members,omitempty"`
	Messages []Message          `bson:"messages"`
}
