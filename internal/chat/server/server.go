package server

import (
	"github.com/Madslick/chat-server/internal/chat/service_assistants"
	"github.com/Madslick/chat-server/pkg"
)

func NewServer(conversationAssistant service_assistants.ConversationAssistant) *ChatroomServer {
	return &ChatroomServer{conversationAssistant: conversationAssistant}
}

type ChatroomServer struct {
	pkg.UnimplementedChatroomServer

	conversationAssistant service_assistants.ConversationAssistant
}
