package main

import (
	//"fmt"
	// "time"
	//"log"
	"github.com/Madslick/chit-chat-go/cmd"
)

func main() {


	cmd.Execute()

	//messages := make(chan string)

	// go func () {

	// 	messages <- "123"

	// 	for i := 0; i < 1000; i++ {
	// 		fmt.Println(i)
	// 	}

	// }()

	// msg := <- messages


	// fmt.Println(msg)

	// go func() {
	// 	for i := 0; i < 1000; i++ {
	// 		fmt.Println(msg)
	// 	}
	// 	messages <- "done"
	// }()
	// done := <- messages
	// fmt.Println(done)
}
