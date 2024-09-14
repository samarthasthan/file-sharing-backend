package main

import (
	"fmt"

	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	grpc_common "github.com/samarthasthan/21BRS1248_Backend/common/grpc"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	zipkinc "github.com/samarthasthan/21BRS1248_Backend/common/zipkin"
	"github.com/samarthasthan/21BRS1248_Backend/services/user/internal/database"
	"github.com/samarthasthan/21BRS1248_Backend/services/user/internal/database/repository"
	grpcin "github.com/samarthasthan/21BRS1248_Backend/services/user/internal/delivery/grpc"
	"google.golang.org/grpc"
)

var (
	USER_GRPC_PORT         string
	USER_DB_PORT           string
	USER_POSTGRES_USER     string
	USER_POSTGRES_PASSWORD string
	USER_POSTGRES_DB       string
	USER_POSTGRES_HOST     string
)

func init() {
	USER_GRPC_PORT = env.GetEnv("USER_GRPC_PORT", "9000")
	USER_DB_PORT = env.GetEnv("USER_DB_PORT", "5432")
	USER_POSTGRES_USER = env.GetEnv("USER_POSTGRES_USER", "root")
	USER_POSTGRES_PASSWORD = env.GetEnv("USER_POSTGRES_PASSWORD", "password")
	USER_POSTGRES_DB = env.GetEnv("USER_POSTGRES_DB", "user-db")
	USER_POSTGRES_HOST = env.GetEnv("USER_POSTGRES_HOST", "localhost")
}

func main() {
	// New Logger
	log := logger.NewLogger("user")
	log.Info("Starting User Service")

	// New Zipkin Tracer
	tracer, _, err := zipkinc.NewTracer("user")
	if err != nil {
		log.Fatalf("Failed to create Zipkin tracer: %v", err)
	}

	// Connect to postgres
	db := database.NewPostgres()
	// Connect to database and Create postgres connection string
	err = db.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", USER_POSTGRES_HOST, USER_DB_PORT, USER_POSTGRES_USER, USER_POSTGRES_PASSWORD, USER_POSTGRES_DB))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}
	// Register Zipkin
	db.RegisterZipkin(tracer)
	defer db.Close()

	// Register repository
	repo := repository.NewRepository(db.Queries)

	service := grpcin.NewUserService(repo, log)

	grpcServer := grpc_common.NewGrpcServer(log, tracer)

	grpcServer.RegisterService(func(s *grpc.Server) {
		proto_go.RegisterUserServiceServer(s, service)
	})

	grpcServer.Run(USER_GRPC_PORT)
}
