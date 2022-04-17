package datastruct

import "github.com/Madslick/chit-chat-go/pkg"

type Client struct {
	Stream   pkg.Chatroom_ConverseServer
	ClientId string
	Name     string
}
