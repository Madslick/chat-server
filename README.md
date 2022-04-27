# chat-server
This is a grpc communication app implemented with a hexagonal architecture (ports and adapters) design. 

* cmd contains the cli commands to run the cli client and the server using the cobra library
* pkg contains generated files like the generated protobuf files
* internal is where the application code all lives. I've separated it into module folders. 1 folder for the time being, chat, with the idea that it could grow and have different modules like auth, etc. 
    * internal/server/ is for the api endpoint handlers to receive requests and hand off processing to the domain layer
    * internal/connectors/ are the application layer that receives the request & provides a layer of separation between the framework layer (internal/server) and the domain layer (internal/services)
    * internal/datastructs/ is the domain layer objects that are used for business logic. 
    * internal/services/ are the services that are responsible for implementing business logic
    * internal/repository/ is for the database. It provides a Repository interface, a default repository to instantiate a repository implementation (mongodb currently) and any repository interface implementations
    * internal/repository/mongostructs has the database structs for mongodb to map database rows to structs in golang to make it easier to work with the data

Run the client with
go run main.go client -s 127.0.0.1:3000

Run the server with
go run main.go server -p 3000
