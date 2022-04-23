package connectors

import (
	"log"

	"github.com/Madslick/chit-chat-go/internal/chat/datastructs"
	"github.com/Madslick/chit-chat-go/internal/chat/pkg"
	"github.com/Madslick/chit-chat-go/internal/chat/services"
)

type ConversationConnector interface {
	CreateConversation(members []*pkg.Client) (*pkg.ConversationResponse, error)
	Converse(stream pkg.Chatroom_ConverseServer) error
}

type conversationConnector struct {
	conversationService services.ConversationService
}

func NewConversationConnector(conversationService services.ConversationService) ConversationConnector {
	return &conversationConnector{
		conversationService: conversationService,
	}
}

func (cc *conversationConnector) CreateConversation(members []*pkg.Client) (*pkg.ConversationResponse, error) {
	// Convert to datastruct
	var chatMembers []*datastructs.Client
	for _, member := range members {
		chatMembers = append(chatMembers, &datastructs.Client{ClientId: member.GetClientId(), Name: member.GetName()})
	}

	// Make call to service function
	conversation, err := cc.conversationService.CreateConversation(chatMembers)
	if err != nil {
		log.Fatal("Error with conversationService.CreateConversation()", err)
	}

	// Convert back to pkg
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

func (cc conversationConnector) Converse(stream pkg.Chatroom_ConverseServer) error {

	return cc.conversationService.Converse(stream)
}
