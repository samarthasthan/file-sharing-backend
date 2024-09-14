package grpc_common

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	address string
}

func init() {
	MAX_FILE_SIZE = env.GetEnvInt("MAX_FILE_SIZE", 500*1024*1024)
}

// NewGrpcClient creates a new gRPC client
func NewGrpcClient(address string) *GrpcClient {
	return &GrpcClient{address: address}
}

// Connect connects the gRPC client to the specified address
func (g *GrpcClient) Connect() error {
	// Use `grpc.WithInsecure()` if you don't have TLS, but it's better to use `grpc.WithTransportCredentials()` with proper TLS credentials in production.
	conn, err := grpc.Dial(g.address,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_FILE_SIZE), grpc.MaxCallSendMsgSize(MAX_FILE_SIZE)))
	if err != nil {
		return err
	}
	g.conn = conn
	return nil
}

// GetConnection returns the gRPC connection
func (g *GrpcClient) GetConnection() *grpc.ClientConn {
	return g.conn
}

// Close closes the gRPC client connection
func (g *GrpcClient) Close() error {
	if g.conn != nil {
		return g.conn.Close()
	}
	return nil
}
