package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Madslick/chit-chat-go/internal/chat/server"
	"github.com/Madslick/chit-chat-go/internal/chat/service_assistants"
	"github.com/Madslick/chit-chat-go/internal/chat/services"
	"github.com/Madslick/chit-chat-go/pkg"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 3000, "Specify the port to serve on")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Instantiate Service Assistants
	conversationAssistant := service_assistants.NewConversationAssistant(services.Conversation())

	s := grpc.NewServer()
	newChatroom := server.NewServer(conversationAssistant)
	pkg.RegisterChatroomServer(s, newChatroom)
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
