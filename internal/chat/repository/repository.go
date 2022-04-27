package repository

import (
	"github.com/Madslick/chit-chat-go/internal/chat/datastructs"
	"github.com/Madslick/chit-chat-go/internal/shared/db"
)

type Repository interface {
	init()
	CreateConversation(members []*datastructs.Client) (string, error)
	CreateMessage(convId string, message datastructs.Message) (bool, error)
	GetConversationByMemberIds(memberIds []string) (datastructs.Conversation, error)
}

func NewRepository(connection db.DbConnection) Repository {
	repo := &mongoRepository{conn: connection}
	repo.init()
	return repo
}
