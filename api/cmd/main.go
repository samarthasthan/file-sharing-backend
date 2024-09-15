package main

import (
	"log"

	"github.com/samarthasthan/21BRS1248_Backend/api/internal/handler"
	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	grpc_common "github.com/samarthasthan/21BRS1248_Backend/common/grpc"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
)

var (
	userClient        proto_go.UserServiceClient
	fileClient        proto_go.FileServiceClient
	err               error
	API_PORT          string
	USER_GRPC_PORT    string
	USER_GRPC_HOST    string
	STORAGE_GRPC_PORT string
	STORAGE_GRPC_HOST string
)

func init() {
	API_PORT = env.GetEnv("API_PORT", "1248")
	USER_GRPC_PORT = env.GetEnv("USER_GRPC_PORT", "9000")
	USER_GRPC_HOST = env.GetEnv("USER_GRPC_HOST", "localhost")
	STORAGE_GRPC_PORT = env.GetEnv("STORAGE_GRPC_PORT", "9002")
	STORAGE_GRPC_HOST = env.GetEnv("STORAGE_GRPC_HOST", "localhost")

	// Initialize the gRPC client
	us := grpc_common.NewGrpcClient(USER_GRPC_HOST + ":" + USER_GRPC_PORT)
	err = us.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	// Create the UserService gRPC client
	userClient = proto_go.NewUserServiceClient(us.GetConnection())

	// Initialize the gRPC client
	fs := grpc_common.NewGrpcClient(STORAGE_GRPC_HOST + ":" + STORAGE_GRPC_PORT)
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
	err := f.Start(API_PORT)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
