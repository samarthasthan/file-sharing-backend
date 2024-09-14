package main

import (
	"log"

	"github.com/samarthasthan/21BRS1248_Backend/api/internal/handler"
	grpc_common "github.com/samarthasthan/21BRS1248_Backend/common/grpc"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
)

var (
	userClient proto_go.UserServiceClient
	fileClient proto_go.FileServiceClient
	err        error
)

func init() {
	// Initialize the gRPC client
	us := grpc_common.NewGrpcClient("localhost:9000")
	err = us.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	// Create the UserService gRPC client
	userClient = proto_go.NewUserServiceClient(us.GetConnection())

	// Initialize the gRPC client
	fs := grpc_common.NewGrpcClient("localhost:9002")
	err = fs.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	fileClient = proto_go.NewFileServiceClient(fs.GetConnection())
}

func main() {
	// Create a Fiber handler with the gRPC client
	f := handler.NewFiberHandler(userClient, fileClient)

	// Register routes
	f.Handle()

	// Start the Fiber app on port 1248
	err := f.Start("1248")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
