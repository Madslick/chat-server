package repository

import (
	"github.com/Madslick/chit-chat-go/internal/chat/datastruct"
	"github.com/Madslick/chit-chat-go/shared/db"
)

type Repository interface {
	CreateConversation(members []*datastruct.Client)
	CreateMessage(message *datastruct.Message)
}

func NewRepository(connection *db.DbConnection) Repository {
	return &mongoRepository{conn: connection}
}
