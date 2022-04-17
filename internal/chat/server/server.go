package server

import (
	"github.com/Madslick/chit-chat-go/internal/chat/service_assistants"
	"github.com/Madslick/chit-chat-go/pkg"
)

func NewServer(conversationAssistant service_assistants.ConversationAssistant) *ChatroomServer {
	return &ChatroomServer{conversationAssistant: conversationAssistant}
}

type ChatroomServer struct {
	pkg.UnimplementedChatroomServer

	conversationAssistant service_assistants.ConversationAssistant
}
