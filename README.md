# chat-server

protoc -I=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative  --go_out=. --go-grpc_out=. chat.proto