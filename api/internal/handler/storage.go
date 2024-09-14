package handler

import (
	"fmt"
	"io/ioutil"
	"os"

	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
)

func (f *FiberHandler) handleFileUpload(c *fiber.Ctx) error {
	// Retrieve the email from Fiber context
	email := c.Locals("email").(string)
	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email is required")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	// Read file data
	fileData, err := readMultipartFile(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read file")
	}

	// Call gRPC service to upload the file
	grpcResponse, err := f.fileClient.UploadFile(c.Context(), &proto_go.UploadFileRequest{
		Email:    email,
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

// handleGetFile handles file download by file ID
func (f *FiberHandler) handleGetFile(c *fiber.Ctx) error {
	// Get the file ID from the URL parameters
	fileID := c.Params("file_id")
	if fileID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File ID is required")
	}

	// Call the gRPC service to get the file metadata (including its storage location)
	grpcResponse, err := f.fileClient.GetFileMetadata(c.Context(), &proto_go.FileMetadataRequest{
		FileId: fileID,
	})
	if err != nil {
		// If there's an error, return it
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to retrieve file metadata: %v", err))
	}

	if grpcResponse.IsProcessed {
		err = c.Redirect(fmt.Sprintf("http://3.7.73.40:13000%s", grpcResponse.StorageLocation))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to serve file: %v", err))
		}
		return nil
	}

	// Construct the file path from the storage location
	filePath := "../../.data" + grpcResponse.StorageLocation

	// Check if the file exists at the given path
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	// Serve the file for download
	err = c.Download(filePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to serve file: %v", err))
	}

	return nil
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
