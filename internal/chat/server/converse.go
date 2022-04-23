package server

import (
	"log"

	"github.com/Madslick/chit-chat-go/internal/chat/pkg"
)

func (cs *ChatroomServer) Converse(stream pkg.Chatroom_ConverseServer) error {
	log.Println("Converse Stream requested")
	return cs.conversationAssistant.Converse(stream)
}
