package mongostructs

type Client struct {
	Id   string `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}
