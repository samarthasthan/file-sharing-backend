package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
)

func (h *FiberHandler) register(c *fiber.Ctx) error {
	var req *proto_go.RegisterRequest

	// Parse the request body into the req struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Use the req struct to call the gRPC client
	res, err := h.userClient.Register(c.Context(), &proto_go.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})

	if err != nil {
		// Set the status code to 500 and return the error message
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *FiberHandler) login(c *fiber.Ctx) error {
	var req *proto_go.LoginRequest

	// Parse the request body into the req struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Use the req struct to call the gRPC client
	res, err := h.userClient.Login(c.Context(), &proto_go.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// Set the status code to 500 and return the error message
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
