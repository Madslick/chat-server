package server

import (
	"context"
	"log"

	"github.com/Madslick/chit-chat-go/internal/chat/pkg"
)

func (cs *ChatroomServer) CreateConversation(ctx context.Context, in *pkg.ConversationRequest) (*pkg.ConversationResponse, error) {

	log.Println("Create Conversation Requested")
	conversationResponse, err := cs.conversationAssistant.CreateConversation(in.GetMembers())

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return conversationResponse, nil
}
