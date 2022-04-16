package services

import (
	"io"
	"log"
	"sync"

	"github.com/Madslick/chat-server/internal/chat/datastruct"
	"github.com/Madslick/chat-server/pkg"
	"github.com/google/uuid"
)

type ConversationService interface {
	CreateConversation(members []*datastruct.Client) (*datastruct.Conversation, error)
	Converse(conversationStream pkg.Chatroom_ConverseServer) error
	Broadcast(from *pkg.Client, conversation *pkg.Conversation, event *pkg.ChatEvent)
}

type conversationService struct {
	clients       sync.Map
	conversations map[string]*datastruct.Conversation
}

var conversationOnce sync.Once
var conversationInstance ConversationService

func Conversation() ConversationService {
	conversationOnce.Do(func() { // <-- atomic, does not allow repeating

		conversationInstance = &conversationService{clients: sync.Map{}, conversations: make(map[string]*datastruct.Conversation)}
	})

	return conversationInstance
}

func (cs *conversationService) CreateConversation(members []*datastruct.Client) (*datastruct.Conversation, error) {

	for _, conversation := range cs.conversations {

		if len(members) != len(conversation.Members) {
			continue
		}

		conversationFound := true
		for _, member := range members {
			memberFound := false
			for _, convMember := range conversation.Members {
				if member.ClientId == convMember.ClientId {
					memberFound = true
				}
			}

			if !memberFound {
				conversationFound = false
			}
		}

		if conversationFound {
			return conversation, nil
		}
	}

	var normalizedMembers []datastruct.Client
	for _, member := range members {
		cs.clients.Range(func(k interface{}, v interface{}) bool {
			clientStream := v.(*datastruct.Client)
			if clientStream.Name == member.Name {
				normalizedMembers = append(normalizedMembers, datastruct.Client{Name: member.Name, ClientId: clientStream.ClientId})
			}
			return true
		})
	}
	conversation := datastruct.Conversation{
		Id:      uuid.New().String(),
		Members: normalizedMembers,
	}

	cs.conversations[conversation.Id] = &conversation
	log.Println("New Conversation saved with Id: ", conversation.Id)
	return &conversation, nil
}

func (cs *conversationService) Converse(stream pkg.Chatroom_ConverseServer) error {

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
			cs.clients.Store(clientId, &datastruct.Client{
				Stream:   stream,
				ClientId: clientId,
				Name:     name,
			})
			from := &pkg.Client{
				ClientId: clientId,
				Name:     name,
			}
			log.Println(name, " logged in with client id: ", from.ClientId)
			cs.Broadcast(from, nil, &pkg.ChatEvent{
				Command: &pkg.ChatEvent_Login{
					Login: from,
				},
			})
		} else if message := in.GetMessage(); message != nil {
			// Get this message broadcasted out
			from := message.GetFrom()
			conversation := message.GetConversation()

			cs.Broadcast(from, conversation, &pkg.ChatEvent{
				Command: &pkg.ChatEvent_Message{
					Message: &pkg.Message{
						From:    from,
						Content: message.GetContent(),
					},
				},
			})
		}
	}
}

func (cs *conversationService) Broadcast(from *pkg.Client, conversation *pkg.Conversation, event *pkg.ChatEvent) {

	if conversation == nil {
		if client, ok := cs.clients.Load(from.ClientId); ok {
			client.(*datastruct.Client).Stream.Send(event)
		} else {
			log.Fatal("Unable to find client for login event by Id: ", from.ClientId)
		}
		return
	}

	for _, to := range conversation.Members {
		if from.ClientId == to.ClientId {
			continue
		}
		if client, ok := cs.clients.Load(to.ClientId); ok {
			client.(*datastruct.Client).Stream.Send(event)
		} else {
			log.Fatal("Unable to find client for message event by Id: ", to.ClientId)
		}

	}
}
