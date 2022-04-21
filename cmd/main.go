package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/Madslick/chit-chat-go/internal/chat/server"
	"github.com/Madslick/chit-chat-go/internal/chat/service_assistants"
	"github.com/Madslick/chit-chat-go/internal/chat/services"
	"github.com/Madslick/chit-chat-go/pkg"
	"github.com/Madslick/chit-chat-go/shared/db"
)

func main() {
	// Receive CLI Arguments
	var port int
	flag.IntVar(&port, "p", 3000, "Specify the port to serve on")
	flag.Parse()

	// Open a Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		os.Exit(1)
	}
	log.Printf("Server listening at %v", lis.Addr())

	// Start gRPC Server
	s := grpc.NewServer()

	// Setup Database
	mongo_client := db.SetupDb()
	if err != nil {
		os.Exit(1)
	}

	// Setup Services
	conversationService := services.Conversation(mongo_client)

	// Setup Service Assistants
	conversationAssistant := service_assistants.NewConversationAssistant(conversationService)

	// Start Chat App
	newChatroom := server.NewServer(conversationAssistant)
	pkg.RegisterChatroomServer(s, newChatroom)

	// Serve the App now
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
