package handler

import (
	"context"
	"io/ioutil"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
)

// Handle file upload by sending a request to the gRPC server
func (f *FiberHandler) handleFileUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}
	// Read file data
	fileData, err := readMultipartFile(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read file 1")
	}

	// Call gRPC service to upload the file
	grpcResponse, err := f.fileClient.UploadFile(context.Background(), &proto_go.UploadFileRequest{
		UserId:   uuid.New().String(), // Assuming user_id is set from JWT
		FileName: file.Filename,
		FileData: fileData,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": grpcResponse.Message,
		"file_id": grpcResponse.FileId,
	})
}

// Handle retrieving file metadata
func (f *FiberHandler) handleGetFileMetadata(c *fiber.Ctx) error {
	fileID := c.Params("file_id")

	grpcResponse, err := f.fileClient.GetFileMetadata(context.Background(), &proto_go.FileMetadataRequest{
		FileId: fileID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get file metadata")
	}

	return c.JSON(grpcResponse)
}

// Helper function to read file from multipart form
func readMultipartFile(file *multipart.FileHeader) ([]byte, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
