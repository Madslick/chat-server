package server

import (
	"context"
	"log"

	"github.com/google/uuid"

	pb "github.com/Madslick/chat-server/pkg"
)

func (cs *ChatroomServer) CreateConversation(ctx context.Context, in *pb.ConversationRequest) (*pb.ConversationResponse, error) {
	members := in.GetMembers()
	var normalizedMembers []*pb.Client
	for _, member := range members {
		cs.clients.Range(func(k interface{}, v interface{}) bool {
			clientStream := v.(*ClientStream)
			if clientStream.Name == member.GetName() {
				normalizedMembers = append(normalizedMembers, &pb.Client{Name: member.GetName(), ClientId: clientStream.ClientId})
			}
			return true
		})
	}
	conversation := pb.Conversation{
		Id:      uuid.New().String(),
		Members: normalizedMembers,
	}

	cs.conversations[conversation.Id] = &conversation
	log.Println("New Conversation saved with Id: ", conversation.Id)

	conversationResponse := pb.ConversationResponse{
		Id:      conversation.Id,
		Members: conversation.Members,
	}
	return &conversationResponse, nil
}
