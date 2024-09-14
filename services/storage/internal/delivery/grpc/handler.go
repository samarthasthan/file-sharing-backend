package grpc

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"os"

	"github.com/google/uuid"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		Expiresat:       sql.NullTime{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file metadata: %v", err)
	}

	return &proto_go.UploadFileResponse{
		FileId:  uuid,
		Message: "File uploaded successfully",
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
