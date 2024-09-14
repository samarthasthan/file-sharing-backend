package api

import (
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/openzipkin/zipkin-go"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	zipkinc "github.com/samarthasthan/21BRS1248_Backend/common/zipkin"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	Start() error
	Handle() error
}

type FiberHandler struct {
	*fiber.App
}

func NewFiberHandler() *FiberHandler {
	// New Tracer
	tracer, _, err := zipkinc.NewTracer("api-gateway")
	if err != nil {
		panic(err)
	}
	// New Logger
	logger := logger.NewLogger("api-gateway")
	logger.Info("Starting API Gateway")
	// New Prometheus
	prometheus := fiberprometheus.New("api-gateway")
	app := fiber.New()
	prometheus.RegisterAt(app, "/metrics")
	app.Use(func(c *fiber.Ctx) error {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"ip":     c.IP(),
		}).Info("Request received")
		return c.Next()
	})
	app.Use(prometheus.Middleware)
	app.Use(zipkinMiddleware(tracer))
	return &FiberHandler{app}
}

func (f *FiberHandler) Start(port string) error {
	return f.Listen(fmt.Sprintf(":%s", port))
}

func (f *FiberHandler) Handle() error {
	f.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	return nil
}

func zipkinMiddleware(tracer *zipkin.Tracer) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		span, _ := tracer.StartSpanFromContext(c.Context(), "handle-request")
		defer span.Finish()
		return c.Next()
	}
}
