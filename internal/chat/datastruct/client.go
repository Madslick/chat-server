package datastruct

import "github.com/Madslick/chat-server/pkg"

type Client struct {
	Stream   pkg.Chatroom_ConverseServer
	ClientId string
	Name     string
}
