package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/samarthasthan/21BRS1248_Backend/common/env"
)

var (
	MINIO_PORT            string
	MINIO_ROOT_USER       string
	MINIO_ROOT_PASSWORD   string
	MINIO_HOST            string
	MINIO_DEFAULT_BUCKETS string
)

func init() {
	MINIO_PORT = env.GetEnv("MINIO_PORT", "13000")
	MINIO_ROOT_USER = env.GetEnv("MINIO_ROOT_USER", "root")
	MINIO_ROOT_PASSWORD = env.GetEnv("MINIO_ROOT_PASSWORD", "password")
	MINIO_HOST = env.GetEnv("MINIO_HOST", "localhost")
	MINIO_DEFAULT_BUCKETS = env.GetEnv("MINIO_DEFAULT_BUCKETS", "uploads")
}

func main() {
	// Configure to use MinIO Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(MINIO_ROOT_USER, MINIO_ROOT_PASSWORD, ""),
		Endpoint:         aws.String(fmt.Sprintf("http://%s:%s", MINIO_HOST, MINIO_PORT)),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	// Create a new session
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create S3 client
	s3Client := s3.New(newSession)

	// Create bucket if it doesn't exist
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(MINIO_DEFAULT_BUCKETS),
	})
	if err != nil {
		if !isErrorBucketAlreadyExists(err) {
			log.Printf("Failed to create bucket: %v", err)
		}
	} else {
		log.Printf("Successfully created bucket: %s", MINIO_DEFAULT_BUCKETS)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(newSession)

	// Upload files from .data/uploads/ directory
	uploadDir := "../../../.data/uploads/"
	err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %v", path, err)
			}
			defer file.Close()

			key := filepath.Base(path)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(MINIO_DEFAULT_BUCKETS),
				Key:    aws.String(key),
				Body:   file,
				ACL:    aws.String("public-read"),
			})
			if err != nil {
				return fmt.Errorf("failed to upload file %s: %v", path, err)
			}

			publicURL := fmt.Sprintf("http://%s:%s/%s/%s", MINIO_HOST, MINIO_PORT, MINIO_DEFAULT_BUCKETS, key)
			log.Printf("Uploaded file: %s", key)
			log.Printf("Public URL: %s", publicURL)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error uploading files: %v", err)
	}

	log.Println("All files uploaded successfully")
}

func isErrorBucketAlreadyExists(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeBucketAlreadyExists:
			return true
		case s3.ErrCodeBucketAlreadyOwnedByYou:
			return true
		default:
			return false
		}
	}
	return false
}
