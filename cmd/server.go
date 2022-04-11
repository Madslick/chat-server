package cmd

// import (
// 	"log"
// 	"net"

// 	pb "github.com/Madslick/chat-server/chat/protos"
// 	"github.com/Madslick/chat-server/chat/server"
// 	"google.golang.org/grpc"
// 	"github.com/spf13/cobra"
// 	"google.golang.org/grpc"
// )

// var serverCmdPort int

// // serverCmd represents the server command
// var serverCmd = &cobra.Command{
// 	Use:   "server",
// 	Short: "Start up chatroom server",
// 	Run: func(cmd *cobra.Command, args[]string) {
// 		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
// 		if err != nil {
// 			log.Fatalf("Failed to listen: %v", err)
// 		}
	
// 		s := grpc.NewServer()
		
// 		// members := []* pb.Client{
// 		// 	&pb.Client{ClientId: "1234", Name: "Joey"},
// 		// 	&pb.Client{ClientId: "5678", Name: "Kira"},
// 		// }
// 		newServer := server.NewServer()
// 		//conversation := newServer.CreateConversation(members)
	
// 		//log.Println("Conversation Id: ", conversation.Id, "- Members: ", len(conversation.Members), " 1st member name: ", conversation.Members[0].Name)
	
// 		pb.RegisterChatroomServer(s, newServer)
// 		log.Printf("Server listening at %v", lis.Addr())
	
// 		if err := s.Serve(lis); err != nil {
// 			log.Fatalf("failed to serve: %v", err)
// 		}
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(serverCmd)
// 	serverCmd.PersistentFlags().IntVarP(&serverCmdPort, "port", "p", 3000, "Port to listen")
// }
	