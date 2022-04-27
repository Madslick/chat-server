package datastructs

type Conversation struct {
	Id       string
	Members  []*Client
	Messages []*Message
}
