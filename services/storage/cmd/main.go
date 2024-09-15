package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	grpc_common "github.com/samarthasthan/21BRS1248_Backend/common/grpc"
	"github.com/samarthasthan/21BRS1248_Backend/common/kafka"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/models"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	zipkinc "github.com/samarthasthan/21BRS1248_Backend/common/zipkin"
	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database"
	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database/repository"
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
	TEMP_PATH                 string
	KAFKA_PORT                string
	KAFKA_HOST                string
	REDIS_HOST                string
	REDIS_PORT                string
)

func init() {
	STORAGE_GRPC_PORT = env.GetEnv("STORAGE_GRPC_PORT", "9002")
	STORAGE_DB_PORT = env.GetEnv("STORAGE_DB_PORT", "5432")
	STORAGE_POSTGRES_STORAGE = env.GetEnv("STORAGE_POSTGRES_STORAGE", "root")
	STORAGE_POSTGRES_PASSWORD = env.GetEnv("STORAGE_POSTGRES_PASSWORD", "password")
	STORAGE_POSTGRES_DB = env.GetEnv("STORAGE_POSTGRES_DB", "postgres")
	STORAGE_POSTGRES_HOST = env.GetEnv("STORAGE_POSTGRES_HOST", "localhost")
	TEMP_PATH = env.GetEnv("TEMP_PATH", "/tmp/uploads")
	KAFKA_PORT = env.GetEnv("KAFKA_PORT", "9092")
	KAFKA_HOST = env.GetEnv("KAFKA_HOST", "localhost")
	REDIS_HOST = env.GetEnv("REDIS_HOST", "localhost")
	REDIS_PORT = env.GetEnv("REDIS_PORT", "6379")
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

	// Initialising Kafka Producer
	p := kafka.NewKafkaProducer(KAFKA_HOST, KAFKA_PORT)

	// Initialising Kafka Consumer
	c := kafka.NewKafkaConsumer(KAFKA_HOST, KAFKA_PORT)
	c.Subscribe([]string{"file-process-out"})

	// Connect to postgres
	db := database.NewPostgres()
	// Connect to database and Create postgres connection string
	err = db.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", STORAGE_POSTGRES_HOST, STORAGE_DB_PORT, STORAGE_POSTGRES_STORAGE, STORAGE_POSTGRES_PASSWORD, STORAGE_POSTGRES_DB))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}
	// Register Zipkin
	db.RegisterZipkin(tracer)
	defer db.Close()

	// Connect to Redis
	rd := database.NewRedis()
	err = rd.Connect(fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT))
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
		panic(err)
	}

	// Register repository
	repo := repository.NewRepository(db.Queries, rd.Client)

	service := grpcin.NewStorageService(log, repo, p, c)

	grpcServer := grpc_common.NewGrpcServer(log, tracer)

	grpcServer.RegisterService(func(s *grpc.Server) {
		proto_go.RegisterFileServiceServer(s, service)
	})

	go func() {
		grpcServer.Run(STORAGE_GRPC_PORT)
	}()

	go func() {
		for {
			msg, err := c.ReadMessage(1 * time.Second)
			if err != nil {
				continue
			} else {
				// Message to models.FileProcess
				fileProcess := &models.FileProcess{}
				if err := json.Unmarshal(msg.Value, fileProcess); err != nil {
					log.Fatalf("Failed to unmarshal message: %v", err)
				}
				// Delete file from local storage
				err = os.Remove(fmt.Sprintf("../../../.data%s", fileProcess.Path))
				if err != nil {
					log.Fatalf("Failed to delete file: %v", err)
				}

				// Update file status in database
				err = repo.MarkFileAsProcessed(context.Background(), fileProcess.ID)
				if err != nil {
					log.Fatalf("Failed to update file status: %v", err)
				}

				// Create a new Mail struct
				mail := &models.Mail{
					To:      fileProcess.Email,
					Subject: fmt.Sprintf("File Processed: %s", fileProcess.ID),
					Body:    fmt.Sprintf("File %s has been processed successfully, public url is http://3.7.73.40:1248/share/%s", fileProcess.ID, fileProcess.ID),
				}

				// Produce a message to the mail topic
				err = p.ProduceMsg(context.Background(), "mail", mail)
				if err != nil {
					log.Fatalf("Failed to produce message: %v", err)
				}

			}
		}
	}()

	go func() {
		// Delete files older than 5 minutes
		for {
			// err := repo.DeleteFile(context.Background())
			time.Sleep(5 * time.Minute)
		}
	}()
}
