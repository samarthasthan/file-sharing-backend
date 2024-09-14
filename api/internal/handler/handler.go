package handler

import (
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
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
}

// NewFiberHandler creates a new Fiber handler with gRPC and Zipkin
func NewFiberHandler(us proto_go.UserServiceClient) *FiberHandler {
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
	app := fiber.New()

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
	return &FiberHandler{app: app, userClient: us}
}

// Start starts the Fiber server on the specified port
func (f *FiberHandler) Start(port string) error {
	return f.app.Listen(fmt.Sprintf(":%s", port))
}

// Handle registers the routes for the Fiber app
func (f *FiberHandler) Handle() {
	f.handleUserRoutes()
}

// Handle user routes
func (f *FiberHandler) handleUserRoutes() {
	// Register user routes
	user := f.app.Group("/user")
	user.Post("/login", f.login)
	user.Post("/register", f.register)
}

// zipkinMiddleware adds Zipkin tracing to the Fiber middleware
func zipkinMiddleware(tracer *zipkin.Tracer) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		span, _ := tracer.StartSpanFromContext(c.Context(), "handle-request")
		defer span.Finish()

		return c.Next()
	}
}
