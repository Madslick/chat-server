package mongostructs

type Account struct {
	Id       string `bson:"_id,omitempty"`
	Email    string `bson:"email,omitempty"`
	Password string `bson:"password,omitempty"`
	Phone    string `bson:"phone,omitempty"`
	First    string `bson:"first,omitempty"`
	Last     string `bson:"last,omitempty"`
}
