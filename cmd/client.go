package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Madslick/chat-server/pkg"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var clientCmdHost string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start up chatroom client",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cc, err := grpc.DialContext(ctx, clientCmdHost, grpc.WithInsecure())
		handleInitError(err, "connect")
		defer cc.Close()

		client := pkg.NewChatroomClient(cc)
		stream, err := client.Converse(context.Background())
		handleInitError(err, "client converse")
		defer stream.CloseSend()

		waitc := make(chan struct{})

		// Take Login Input
		log.Println("Who are you?")
		mainScanner := bufio.NewScanner(os.Stdin)
		mainScanner.Scan()
		name := strings.TrimSpace(mainScanner.Text())

		// Send Login over to server
		chatClient := &pkg.Client{Name: name}
		loginEvent := pkg.ChatEvent{
			Command: &pkg.ChatEvent_Login{
				Login: chatClient,
			},
		}
		sendErr := stream.Send(&loginEvent)
		if sendErr != nil {
			log.Fatalf("Failed to send message to server: %v", err)
		}

		// Receive Login Response
		loginResponse, _ := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to login to server")
		}

		if login := loginResponse.GetLogin(); login != nil {
			chatClient.ClientId = login.GetClientId()
		}

		log.Println("Who do you wanna call?")
		mainScanner.Scan()
		memberName := strings.TrimSpace(mainScanner.Text())
		conversation := pkg.Conversation{
			Members: []*pkg.Client{chatClient, &pkg.Client{Name: memberName}},
		}

		// Receiving message from server
		go func() {
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					close(waitc)
					return
				}
				if err != nil {
					log.Fatalf("Failed to receive a note : %v", err)
				}

				if login := in.GetLogin(); login != nil {
					log.Println(login.GetName(), "logged in")
				} else if message := in.GetMessage(); message != nil {
					log.Println(message.GetFrom().GetName(), ":", message.GetContent())
				}
			}
		}()

		// Reading message from stdin and send to server
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())
				if text == "" {
					continue
				}
				// event := nil
				// if event == nil {
				// 	continue
				// }
				message := pkg.Message{}
				message.Conversation = &conversation
				message.From = chatClient
				message.Content = text

				err := stream.Send(&pkg.ChatEvent{
					Command: &pkg.ChatEvent_Message{Message: &message},
				})
				if err != nil {
					log.Fatalf("Failed to send message to server: %v", err)
				}
			}
		}()
		log.Println("client started")

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	loop:
		for {
			select {
			case <-waitc:
				break loop
			case <-quit:
				break loop
			}
		}
		log.Println("client exited")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().StringVarP(&clientCmdHost, "host", "s", "127.0.0.1", "Server host address")
}
