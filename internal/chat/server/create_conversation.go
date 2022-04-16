package server

import (
	"context"
	"log"

	"github.com/Madslick/chat-server/pkg"
)

func (cs *ChatroomServer) CreateConversation(ctx context.Context, in *pkg.ConversationRequest) (*pkg.ConversationResponse, error) {

	conversationResponse, err := cs.conversationAssistant.CreateConversation(in.GetMembers())

	if err != nil {
		log.Fatal(err)
	}

	return conversationResponse, nil
}
