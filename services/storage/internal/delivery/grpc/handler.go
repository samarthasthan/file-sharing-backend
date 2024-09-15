package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"os"

	"github.com/google/uuid"
	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	"github.com/samarthasthan/21BRS1248_Backend/common/models"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	HOST string
)

func init() {
	HOST = env.GetEnv("HOST", "localhost:1248")

}

// File storage path (local disk)
const fileStoragePath = "../../../.data/uploads/"

// UploadFile handles file upload and saves metadata in PostgreSQL
func (s *StorageService) UploadFile(ctx context.Context, req *proto_go.UploadFileRequest) (*proto_go.UploadFileResponse, error) {

	uuid := uuid.New().String()
	filePath := fileStoragePath + uuid + "_" + req.GetFileName()

	os.MkdirAll(fileStoragePath, os.ModePerm)
	err := os.WriteFile(filePath, req.GetFileData(), 0644)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file: %v", err)
	}

	fileType := http.DetectContentType(req.GetFileData())
	fileSize := len(req.GetFileData())

	err = s.repo.UploadFile(ctx, &sqlc.UploadFileByEmailParams{
		Fileid:          uuid,
		Email:           req.Email,
		Filename:        req.GetFileName(),
		Filetype:        fileType,
		Filesize:        int64(fileSize),
		Storagelocation: fmt.Sprintf("/uploads/%s", uuid+"_"+req.GetFileName()),
		Uploaddate:      time.Now(),
		Expiresat:       sql.NullTime{Time: time.Now().Add(time.Minute * 5), Valid: true},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file metadata: %v", err)
	}

	go func() {
		// Publish message to Kafka
		s.p.ProduceMsg(context.Background(), "file-process-in", &models.FileProcess{
			ID:    uuid,
			Path:  fmt.Sprintf("/uploads/%s", uuid+"_"+req.GetFileName()),
			Email: req.Email,
		})
	}()

	return &proto_go.UploadFileResponse{
		FileId:    uuid,
		Message:   "File uploaded successfully",
		PublicUrl: fmt.Sprintf("%s/share/%s", HOST, uuid),
	}, nil
}

// GetFile handles file download
func (s *StorageService) GetFileMetadata(ctx context.Context, req *proto_go.FileMetadataRequest) (*proto_go.FileMetadataResponse, error) {
	file, err := s.repo.GetFile(ctx, req.GetFileId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	}

	return &proto_go.FileMetadataResponse{
		IsProcessed:     file.Isprocessed.Bool,
		StorageLocation: file.Storagelocation,
	}, nil
}

// GetFilesByUser handles the retrieval of files by user with caching
func (s *StorageService) GetFilesByUser(ctx context.Context, req *proto_go.FilesByUserRequest) (*proto_go.FilesByUserResponse, error) {
	// Redis cache key based on user email
	cacheKey := fmt.Sprintf("user:%s:files", req.GetEmail())

	// Try to get the data from Redis
	cachedFiles, err := s.repo.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedFiles != "" {
		s.log.Info("Cache hit")
		// If cache hit, unmarshal the cached data
		var filesResp []*proto_go.File
		if err := json.Unmarshal([]byte(cachedFiles), &filesResp); err == nil {
			// Return the cached response
			return &proto_go.FilesByUserResponse{
				Files: filesResp,
			}, nil
		}
	}

	// If cache miss, query the database
	files, err := s.repo.FilesByUserRequest(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "files not found: %v", err)
	}

	// Prepare the response
	var filesResp []*proto_go.File
	for _, file := range files {
		filesResp = append(filesResp, &proto_go.File{
			FileId:          file.Fileid,
			FileName:        file.Filename,
			FileSize:        fmt.Sprintf("%d bytes", file.Filesize),
			FileType:        file.Filetype,
			StorageLocation: file.Storagelocation,
			UploadDate:      file.Uploaddate.String(),
			IsProcessed:     file.Isprocessed.Bool,
			ExpiredAt:       file.Expiresat.Time.String(),
			PublicUrl:       fmt.Sprintf("%s/share/%s", HOST, file.Fileid),
		})
	}

	// Marshal the response to store in Redis
	cachedData, err := json.Marshal(filesResp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cache files: %v", err)
	}

	// Store the data in Redis with an expiration time (e.g., 1 hour)
	err = s.repo.Redis.Set(ctx, cacheKey, cachedData, time.Minute*5).Err()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set cache: %v", err)
	}

	// Return the fresh data
	return &proto_go.FilesByUserResponse{
		Files: filesResp,
	}, nil
}
