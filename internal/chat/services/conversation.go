package services

import (
	"io"
	"log"
	"sync"

	"github.com/google/uuid"

	"github.com/Madslick/chit-chat-go/internal/chat/datastruct"
	"github.com/Madslick/chit-chat-go/pkg"
	"github.com/Madslick/chit-chat-go/shared/db"
)

type ConversationService interface {
	CreateConversation(members []*datastruct.Client) (*datastruct.Conversation, error)
	Converse(conversationStream pkg.Chatroom_ConverseServer) error
	Broadcast(from *pkg.Client, conversation *pkg.Conversation, event *pkg.ChatEvent)
}

type conversationService struct {
	mongoClient   db.DbConnection
	clients       sync.Map
	conversations map[string]*datastruct.Conversation
}

var conversationOnce sync.Once
var conversationInstance ConversationService

func Conversation(mongoClient db.DbConnection) ConversationService {
	conversationOnce.Do(func() { // <-- atomic, does not allow repeating

		conversationInstance = &conversationService{
			mongoClient:   mongoClient,
			clients:       sync.Map{},
			conversations: make(map[string]*datastruct.Conversation)}
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
				if member.Name == convMember.Name {
					memberFound = true
				}
			}

			if !memberFound {
				conversationFound = false
			}
		}

		if conversationFound {
			log.Printf("Returning Conversation %s found in cache\n", conversation.Id)
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
	var clientId string
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			cs.clients.Delete(clientId)
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
			cs.Broadcast(from, nil, &pkg.ChatEvent{
				Command: &pkg.ChatEvent_Login{
					Login: from,
				},
			})
			log.Printf("%s logged in with client id: %s\n", name, from.GetClientId())
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
			log.Printf("Message sent from %s in conversation %s with content %s\n", from.GetName(), conversation.GetId(), message.GetContent())
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
