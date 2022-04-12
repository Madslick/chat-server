package cmd

import (
	"log"
	"fmt"
	"net"

	pb "github.com/Madslick/chat-server/chat/protos"
	"github.com/Madslick/chat-server/chat/server"
	"google.golang.org/grpc"
	"github.com/spf13/cobra"
)

var serverCmdPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start up chatroom server",
	Run: func(cmd *cobra.Command, args[]string) {
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
