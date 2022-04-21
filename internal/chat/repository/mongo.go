package repository

import (
	"github.com/Madslick/chit-chat-go/internal/chat/datastruct"
	"github.com/Madslick/chit-chat-go/shared/db"
)

type mongoRepository struct {
	conn *db.DbConnection
}

func (mr *mongoRepository) CreateConversation(members []*datastruct.Client) {

}

func (mr *mongoRepository) CreateMessage(message *datastruct.Message) {

}
