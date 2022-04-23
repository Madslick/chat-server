package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/Madslick/chit-chat-go/internal/shared/db"

	chatConnectors "github.com/Madslick/chit-chat-go/internal/chat/connectors"
	chatPkg "github.com/Madslick/chit-chat-go/internal/chat/pkg"
	chatRepository "github.com/Madslick/chit-chat-go/internal/chat/repository"
	chatServer "github.com/Madslick/chit-chat-go/internal/chat/server"
	chatServices "github.com/Madslick/chit-chat-go/internal/chat/services"

	authConnectors "github.com/Madslick/chit-chat-go/internal/auth/connectors"
	authPkg "github.com/Madslick/chit-chat-go/internal/auth/pkg"
	authRepository "github.com/Madslick/chit-chat-go/internal/auth/repository"
	authServer "github.com/Madslick/chit-chat-go/internal/auth/server"
	authServices "github.com/Madslick/chit-chat-go/internal/auth/services"
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
	dbConnection := db.SetupDb()
	if err != nil {
		os.Exit(1)
	}

	chatRepo := chatRepository.NewRepository(dbConnection)
	authRepo := authRepository.NewRepository(dbConnection)

	// Setup Services
	conversationService := chatServices.NewConversationService(chatRepo)
	accountService := authServices.NewAccountService(authRepo)

	// Setup Server-to-Service Connectors
	conversationConnector := chatConnectors.NewConversationConnector(conversationService)
	accountConnector := authConnectors.NewAccountConnector(accountService)

	// Start Chat App
	newChatroom := chatServer.NewServer(conversationConnector)
	newAuth := authServer.NewServer(accountConnector)

	chatPkg.RegisterChatroomServer(s, newChatroom)
	authPkg.RegisterAuthServer(s, newAuth)

	// Serve the App now
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
