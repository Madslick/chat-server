package mongostructs

type Client struct {
	Id   string `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
}
