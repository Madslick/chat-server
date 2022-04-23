package datastructs

type Message struct {
	from         Client
	content      string
	conversation Conversation
}
