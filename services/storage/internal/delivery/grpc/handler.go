package grpc

import (
	"context"

	"os"

	"github.com/google/uuid"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// File storage path (local disk)
const fileStoragePath = "../../../.data/uploads/"

// UploadFile handles file upload and saves metadata in PostgreSQL
func (s *StorageService) UploadFile(ctx context.Context, req *proto_go.UploadFileRequest) (*proto_go.UploadFileResponse, error) {

	os.MkdirAll(fileStoragePath, os.ModePerm)
	// Save file to local disk
	filePath := fileStoragePath + req.GetFileName()
	err := os.WriteFile(filePath, req.GetFileData(), 0644)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file: %v", err)
	}

	// // Save file metadata to PostgreSQL
	// fileID, err := db.SaveFileMetadata(req.GetUserId(), req.GetFileName(), filePath, len(req.GetFileData()))
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to save file metadata: %v", err)
	// }

	return &proto_go.UploadFileResponse{
		Success: true,
		Message: "File uploaded successfully",
		FileId:  uuid.New().String(),
	}, nil
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
