package grpc

import (
	"context"
	"database/sql"
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
	os.MkdirAll(fileStoragePath, os.ModePerm)
	// Save file to local disk
	filePath := fileStoragePath + uuid + "_" + req.GetFileName()
	err := os.WriteFile(filePath, req.GetFileData(), 0644)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file: %v", err)
	}

	// Save file metadata to PostgreSQL
	var res *proto_go.UploadFileResponse
	err = s.repo.UploadFile(ctx, &sqlc.UploadFileByEmailParams{
		Email:           "samarthasthan27@gmail.com",
		Filename:        req.GetFileName(),
		Filetype:        "image/jpeg",
		Filesize:        123456,
		Storagelocation: filePath,
		Uploaddate:      time.Now(),
		Expiresat:       sql.NullTime{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file metadata: %v", err)
	}
	return res, nil
}

// GetFileMetadata handles retrieving file metadata from PostgreSQL
func (s *StorageService) GetFileMetadata(ctx context.Context, req *proto_go.FileMetadataRequest) (*proto_go.FileMetadataResponse, error) {
	// fileMeta, err := db.GetFileMetadata(req.GetFileId())
	// if err != nil {
	// 	return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	// }

	// return &proto_go.FileMetadataResponse{
	// 	FileName:        fileMeta.FileName,
	// 	UploadDate:      fileMeta.UploadDate.Format(time.RFC3339),
	// 	FileSize:        fileMeta.FileSize,
	// 	StorageLocation: fileMeta.StorageLocation,
	// }, nil
	return nil, nil
}
