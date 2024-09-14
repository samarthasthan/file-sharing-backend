package handler

import (
	"io/ioutil"

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

// handleGetFile
func (f *FiberHandler) handleGetFile(c *fiber.Ctx) error {
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
