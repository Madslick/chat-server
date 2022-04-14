package server

import (
	"sync"

	pb "github.com/Madslick/chat-server/pkg"
)

func NewServer() *ChatroomServer {
	return &ChatroomServer{clients: sync.Map{}, conversations: make(map[string]*pb.Conversation)}
}

type ClientStream struct {
	Stream   pb.Chatroom_ConverseServer
	ClientId string
	Name     string
}

type ChatroomServer struct {
	pb.UnimplementedChatroomServer

	clients       sync.Map
	conversations map[string]*pb.Conversation
}
