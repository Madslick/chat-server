package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/Madslick/chat-server/internal/chat/server"
	pb "github.com/Madslick/chat-server/pkg"
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

		s := grpc.NewServer()

		pb.RegisterChatroomServer(s, server.NewServer())
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
