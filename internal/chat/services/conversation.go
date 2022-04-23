package services

import (
	"io"
	"log"
	"sync"

	"github.com/Madslick/chit-chat-go/internal/chat/datastructs"
	"github.com/Madslick/chit-chat-go/internal/chat/pkg"
	"github.com/Madslick/chit-chat-go/internal/chat/repository"
)

type ConversationService interface {
	CreateConversation(members []*datastructs.Client) (*datastructs.Conversation, error)
	Converse(conversationStream pkg.Chatroom_ConverseServer) error
	Broadcast(from *pkg.Client, conversation *pkg.Conversation, event *pkg.ChatEvent)
}

type conversationService struct {
	repo    repository.Repository
	clients sync.Map
}

var conversationOnce sync.Once
var conversationInstance ConversationService

func NewConversationService(repo repository.Repository) ConversationService {
	conversationOnce.Do(func() { // <-- atomic, does not allow repeating

		conversationInstance = &conversationService{
			repo:    repo,
			clients: sync.Map{},
		}
	})

	return conversationInstance
}

func (cs *conversationService) CreateConversation(members []*datastructs.Client) (*datastructs.Conversation, error) {
	var memberNames []string
	for _, member := range members {
		memberNames = append(memberNames, member.Name)
	}
	conversation, err := cs.repo.GetConversationByMemberNames(memberNames)
	if err != nil {
		log.Fatal(err)
	}
	if conversation.Id != "" {
		log.Println("Found Conversation in database!")
		return &conversation, nil
	}

	id, insertError := cs.repo.CreateConversation(members)
	if insertError != nil {
		log.Fatal("Insert Error calling the repo.CreateConversation from the conversation Service", insertError)
	}

	// CLIENT IS LOOKING AT MEMBERS CLIENTID INSTEAD OF ID
	conversation = datastructs.Conversation{
		Id:      id,
		Members: members,
	}
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
			clientId := login.GetClientId()
			name := login.GetName()
			log.Printf("%s with %s has requested a stream registration", name, clientId)
			cs.clients.Store(clientId, &datastructs.Client{
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
			client.(*datastructs.Client).Stream.Send(event)
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
			client.(*datastructs.Client).Stream.Send(event)
		} else {
			log.Fatalf("Unable to find client for message event by Id: %s", to.ClientId)
		}
	}
}
