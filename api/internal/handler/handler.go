package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/openzipkin/zipkin-go"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	zipkinc "github.com/samarthasthan/21BRS1248_Backend/common/zipkin"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	Start(port string) error
	Handle() error
}

type FiberHandler struct {
	app        *fiber.App
	userClient proto_go.UserServiceClient
	fileClient proto_go.FileServiceClient
}

// NewFiberHandler creates a new Fiber handler with gRPC and Zipkin
func NewFiberHandler(us proto_go.UserServiceClient, fl proto_go.FileServiceClient) *FiberHandler {
	// Initialize Zipkin tracer
	tracer, _, err := zipkinc.NewTracer("api-gateway")
	if err != nil {
		panic(err)
	}

	// Initialize logger
	logger := logger.NewLogger("api-gateway")
	logger.Info("Starting API Gateway")

	// Initialize Prometheus metrics for Fiber
	prometheus := fiberprometheus.New("api-gateway")
	app := fiber.New(
		fiber.Config{
			BodyLimit: 1024 * 1024 * 10, // 10MB limit for file uploads
		},
	)

	// Register Prometheus metrics endpoint
	prometheus.RegisterAt(app, "/metrics")

	// Add request logging middleware
	app.Use(func(c *fiber.Ctx) error {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"ip":     c.IP(),
		}).Info("Request received")
		return c.Next()
	})

	// Add Prometheus and Zipkin middleware
	app.Use(prometheus.Middleware)
	app.Use(zipkinMiddleware(tracer))

	// Return the Fiber handler
	return &FiberHandler{app: app, userClient: us, fileClient: fl}
}

// JWT middleware to protect routes
func jwtMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("secret"), // Replace "secret" with your actual JWT secret
		ErrorHandler: jwtError,
	})
}

// JWT error handler
func jwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Next()
}

// Start starts the Fiber server on the specified port
func (f *FiberHandler) Start(port string) error {
	return f.app.Listen(fmt.Sprintf(":%s", port))
}

// Handle registers the routes for the Fiber app
func (f *FiberHandler) Handle() {
	f.handleUserRoutes()
	f.handleFileRoutes()
}

// Handle user routes
func (f *FiberHandler) handleUserRoutes() {
	// Register user routes
	user := f.app.Group("/user")
	user.Post("/login", f.login)
	user.Post("/register", f.register)
}

// Handle file routes
func (f *FiberHandler) handleFileRoutes() {

	f.app.Get("/share/:file_id", f.handleGetFile)

	// Add auth middleware
	f.app.Use(authMiddleware(f.userClient))
	// Protect routes with JWT middleware
	f.app.Post("/upload", jwtMiddleware(), f.handleFileUpload)
	f.app.Get("/files", jwtMiddleware(), f.handleGetFiles)
}

// zipkinMiddleware adds Zipkin tracing to the Fiber middleware
func zipkinMiddleware(tracer *zipkin.Tracer) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		span, _ := tracer.StartSpanFromContext(c.Context(), "handle-request")
		defer span.Finish()

		return c.Next()
	}
}

func authMiddleware(userClient proto_go.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is missing",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		req := &proto_go.CheckJWTRequest{
			SessionId: token,
		}

		resp, err := userClient.CheckJWT(context.Background(), req)
		if err != nil || !resp.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("email", resp.Email)

		return c.Next()
	}
}
