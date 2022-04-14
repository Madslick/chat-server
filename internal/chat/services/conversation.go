package services

import (
	pb "github.com/Madslick/chat-server/pkg"
	"github.com/Madslick/chat-server/internal/chat/dto"
)

type ConversationService interface {
	CreateConversation(members dto.Client) (*datastruct.Conversation, error)
}
type conversationService struct {}

func (cs *conversationService) CreateConversation(members dto.Client)) ()


