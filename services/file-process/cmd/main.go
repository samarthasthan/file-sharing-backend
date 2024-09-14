package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	"github.com/samarthasthan/21BRS1248_Backend/common/kafka"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/models"
	"github.com/samarthasthan/21BRS1248_Backend/services/file-process/internal/s3"
)

var (
	MINIO_PORT            string
	MINIO_ROOT_USER       string
	MINIO_ROOT_PASSWORD   string
	MINIO_HOST            string
	MINIO_DEFAULT_BUCKETS string
	KAFKA_PORT            string
	KAFKA_HOST            string
)

func init() {
	MINIO_PORT = env.GetEnv("MINIO_PORT", "13000")
	MINIO_ROOT_USER = env.GetEnv("MINIO_ROOT_USER", "root")
	MINIO_ROOT_PASSWORD = env.GetEnv("MINIO_ROOT_PASSWORD", "password")
	MINIO_HOST = env.GetEnv("MINIO_HOST", "localhost")
	MINIO_DEFAULT_BUCKETS = env.GetEnv("MINIO_DEFAULT_BUCKETS", "uploads")
	KAFKA_PORT = env.GetEnv("KAFKA_PORT", "9092")
	KAFKA_HOST = env.GetEnv("KAFKA_HOST", "localhost")
}

func main() {

	// Initialize logger
	log := logger.NewLogger("file-process")

	// Kakfa producer
	p := kafka.NewKafkaProducer(KAFKA_HOST, KAFKA_PORT)

	// Kafka consumer
	c := kafka.NewKafkaConsumer(KAFKA_HOST, KAFKA_PORT)
	c.Subscribe([]string{"file-process-in"})

	// Initialize S3 client
	s3Client, err := s3.NewS3(
		MINIO_ROOT_USER,
		MINIO_ROOT_PASSWORD,
		fmt.Sprintf("http://%s:%s", MINIO_HOST, MINIO_PORT),
		"us-east-1",
		MINIO_DEFAULT_BUCKETS,
	)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	// Create bucket if it doesn't exist
	if err := s3Client.CreateBucket(); err != nil {
		log.Fatalf("Error creating bucket: %v", err)
	}

	for {
		msg, err := c.ReadMessage(1 * time.Second)
		if err != nil {
			continue
		} else {
			log.Infof("Received message: %s", string(msg.Value))
			// Convert message to FileProcess struct
			var fileProcess *models.FileProcess
			if err := json.Unmarshal(msg.Value, &fileProcess); err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}

			// Upload files from .data/uploads/ directory
			uploadDir := fmt.Sprintf("../../../.data%s", fileProcess.Path)
			err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					publicURL, err := s3Client.UploadFile(path)
					if err != nil {
						return fmt.Errorf("failed to upload file %s: %v", path, err)
					}

					log.Printf("Uploaded file: %s", path)
					log.Printf("Public URL: %s", publicURL)
				}
				return nil
			})

			if err != nil {
				log.Fatalf("Error uploading files: %v", err)
			}

			log.Println("All files uploaded successfully")

			p.ProduceMsg(context.Background(), "file-process-out", &fileProcess)
		}

	}
}
