package server

import (
	"github.com/Madslick/chat-server/pkg"
)

func (cs *ChatroomServer) Converse(stream pkg.Chatroom_ConverseServer) error {
	return cs.conversationAssistant.Converse(stream)
}
