package cmd

import (
	"fmt"
	"os"
	"log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chatroom",
	Short: "Chatroom is a simple example for gRPC bidirectional streaming.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleInitError(err error, module string) {
	if err == nil {
		return
	}
	log.Fatalf("init %s failed, err: %s", module, err)
}