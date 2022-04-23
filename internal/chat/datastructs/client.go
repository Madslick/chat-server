package datastructs

import "github.com/Madslick/chit-chat-go/internal/chat/pkg"

type Client struct {
	Stream   pkg.Chatroom_ConverseServer
	ClientId string
	Name     string
}
