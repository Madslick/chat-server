# chat-server

1. cmd contains the cli commands to run the cli client and the server using the cobra library
2. pkg contains generated files like the generated protobuf files
3. internal is where the application code all lives. I've separated it into 1 folder for the time being, chat, with the idea that it could grow and have different modules like auth, etc.

Run the client with
go run main.go client -s 127.0.0.1:3000

Run the server with
go run main.go server -p 3000
