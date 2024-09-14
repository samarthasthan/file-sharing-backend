package main

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	grpc_common "github.com/samarthasthan/21BRS1248_Backend/common/grpc"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	zipkinc "github.com/samarthasthan/21BRS1248_Backend/common/zipkin"
	grpcin "github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/delivery/grpc"
	"google.golang.org/grpc"
)

var (
	STORAGE_GRPC_PORT         string
	STORAGE_DB_PORT           string
	STORAGE_POSTGRES_STORAGE  string
	STORAGE_POSTGRES_PASSWORD string
	STORAGE_POSTGRES_DB       string
	STORAGE_POSTGRES_HOST     string
)

func init() {
	STORAGE_GRPC_PORT = env.GetEnv("STORAGE_GRPC_PORT", "9002")
	STORAGE_DB_PORT = env.GetEnv("STORAGE_DB_PORT", "5432")
	STORAGE_POSTGRES_STORAGE = env.GetEnv("STORAGE_POSTGRES_STORAGE", "root")
	STORAGE_POSTGRES_PASSWORD = env.GetEnv("STORAGE_POSTGRES_PASSWORD", "password")
	STORAGE_POSTGRES_DB = env.GetEnv("STORAGE_POSTGRES_DB", "STORAGE-db")
	STORAGE_POSTGRES_HOST = env.GetEnv("STORAGE_POSTGRES_HOST", "localhost")
}

func main() {
	// New Logger
	log := logger.NewLogger("user")
	log.Info("Starting Storage Service")

	// New Zipkin Tracer
	tracer, _, err := zipkinc.NewTracer("user")
	if err != nil {
		log.Fatalf("Failed to create Zipkin tracer: %v", err)
	}

	// // Connect to postgres
	// db := database.NewPostgres()
	// // Connect to database and Create postgres connection string
	// err = db.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", USER_POSTGRES_HOST, USER_DB_PORT, USER_POSTGRES_USER, USER_POSTGRES_PASSWORD, USER_POSTGRES_DB))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// 	panic(err)
	// }
	// // Register Zipkin
	// db.RegisterZipkin(tracer)
	// defer db.Close()

	// Register repository
	// repo := repository.NewRepository(db.Queries)

	service := grpcin.NewStorageService(log)

	grpcServer := grpc_common.NewGrpcServer(log, tracer)

	grpcServer.RegisterService(func(s *grpc.Server) {
		proto_go.RegisterFileServiceServer(s, service)
	})

	grpcServer.Run(STORAGE_GRPC_PORT)
}
