package server;

import (
	"io"
	"sync"
	"log"

	pb "github.com/Madslick/chat-server/chat/protos"
	"github.com/google/uuid"
)


func NewServer() *ChatroomServer {
	return &ChatroomServer{clients: sync.Map{}, conversations: make(map[string]pb.Conversation)}
}

type ClientStream struct {
	Stream pb.Chatroom_ConverseServer
	ClientId string
	Name string
}

type ChatroomServer struct {
	pb.UnimplementedChatroomServer

	clients sync.Map
	conversations map[string] pb.Conversation
}

func (cs *ChatroomServer) CreateConversation(members []*pb.Client) *pb.Conversation {
	var normalizedMembers []*pb.Client
	for _, member := range members {
		cs.clients.Range(func(k interface{}, v interface{}) bool {
			clientStream := v.(*ClientStream)
			if clientStream.Name == member.GetName() {
				normalizedMembers = append(normalizedMembers, &pb.Client{Name: member.GetName(), ClientId: clientStream.ClientId,},)
			}
			return true
		})
	}
	conversation := pb.Conversation{
		Id: uuid.New().String(),
		Members: normalizedMembers,
	}

	cs.conversations[conversation.Id] = conversation
	log.Println("New Conversation saved with Id: ", conversation.Id)
	return &conversation
}

func (cs *ChatroomServer) Converse(stream pb.Chatroom_ConverseServer) error {
	for {
		in, err := stream.Recv()
		var clientId string

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if login := in.GetLogin(); login != nil {
			clientId = uuid.New().String()
			name := login.GetName()
			cs.clients.Store(clientId, &ClientStream{
				Stream: stream,
				ClientId: clientId,
				Name: name,
			})
			from := &pb.Client{
				ClientId: clientId,
				Name: name,
			}
			log.Println(name, " logged in with client id: ", from.ClientId)
			cs.Broadcast(from, nil, &pb.ChatEvent{
				Command: &pb.ChatEvent_Login{
					Login: from,
				},
			})
		} else if message := in.GetMessage(); message != nil {
			// Get this message broadcasted out
			from := message.GetFrom()
			conversation := message.GetConversation()
			conversationId := conversation.GetId()
			if conversationId == "" {
				conversation = cs.CreateConversation(conversation.GetMembers())
			}

			cs.Broadcast(from, conversation, &pb.ChatEvent{
				Command: &pb.ChatEvent_Message{
					Message: &pb.Message{
						From: from,
						Content: message.GetContent(),
					},
				},
			})
		}
	}
}

func (cs *ChatroomServer) Broadcast(from *pb.Client, conversation *pb.Conversation, event *pb.ChatEvent) {

	if conversation == nil {
		if client, ok := cs.clients.Load(from.GetClientId()); ok {
			client.(*ClientStream).Stream.Send(event)
		} else {
			log.Fatalf("Unable to find client for login event by Id: ", from.GetClientId())
		}
		return
	}

	for _, to := range conversation.Members {
		if from.GetClientId() == to.GetClientId() {
			continue
		}
		if client, ok := cs.clients.Load(to.GetClientId()); ok {
			client.(*ClientStream).Stream.Send(event)
		} else {
			log.Fatalf("Unable to find client for message event by Id: ", to.GetClientId())
		}
		
	}
}

