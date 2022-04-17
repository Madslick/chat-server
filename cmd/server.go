package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/Madslick/chit-chat-go/internal/chat/server"
	"github.com/Madslick/chit-chat-go/internal/chat/service_assistants"
	"github.com/Madslick/chit-chat-go/internal/chat/services"
	"github.com/Madslick/chit-chat-go/pkg"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmdPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start up chatroom server",
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCmdPort))
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
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().IntVarP(&serverCmdPort, "port", "p", 3000, "Port to listen")
}
