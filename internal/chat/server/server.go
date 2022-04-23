package server

import (
	"github.com/Madslick/chit-chat-go/internal/chat/connectors"
	"github.com/Madslick/chit-chat-go/internal/chat/pkg"
)

func NewServer(conversationAssistant connectors.ConversationConnector) *ChatroomServer {
	return &ChatroomServer{conversationAssistant: conversationAssistant}
}

type ChatroomServer struct {
	pkg.UnimplementedChatroomServer

	conversationAssistant connectors.ConversationConnector
}
