package service_assistants

import (
	"log"

	"github.com/Madslick/chat-server/internal/chat/datastruct"
	"github.com/Madslick/chat-server/internal/chat/services"
	"github.com/Madslick/chat-server/pkg"
)

type ConversationAssistant interface {
	CreateConversation(members []*pkg.Client) (*pkg.ConversationResponse, error)
	Converse(stream pkg.Chatroom_ConverseServer) error
}

type conversationAssistant struct {
	conversationService services.ConversationService
}

func NewConversationAssistant(conversationService services.ConversationService) ConversationAssistant {
	return &conversationAssistant{
		conversationService: conversationService,
	}
}

func (ca *conversationAssistant) CreateConversation(members []*pkg.Client) (*pkg.ConversationResponse, error) {
	var chatMembers []*datastruct.Client
	for _, member := range members {
		chatMembers = append(chatMembers, &datastruct.Client{ClientId: member.GetClientId(), Name: member.GetName()})
	}
	conversation, err := ca.conversationService.CreateConversation(chatMembers)
	if err != nil {
		log.Fatal("Error with conversationService.CreateConversation()", err)
	}

	var responseMembers []*pkg.Client
	for _, member := range conversation.Members {
		responseMembers = append(responseMembers, &pkg.Client{Name: member.Name, ClientId: member.ClientId})
	}

	conversationResponse := &pkg.ConversationResponse{
		Id:      conversation.Id,
		Members: responseMembers,
	}

	return conversationResponse, err
}

func (ca conversationAssistant) Converse(stream pkg.Chatroom_ConverseServer) error {

	return ca.conversationService.Converse(stream)
}
