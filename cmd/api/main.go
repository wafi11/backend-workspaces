package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/wafi11/backend-workspaces/pkg/config"
	"github.com/wafi11/backend-workspaces/pkg/middlewares"
	"github.com/wafi11/backend-workspaces/pkg/server"
)

func main() {
	fmt.Println("Application Starting...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := cfg.Database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Database connected successfully!")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(middlewares.CustomLogger(middlewares.LoggerConfig{
		Skip: func(c *fiber.Ctx) bool {
			return c.Path() == "/health" || c.Path() == "/ping"
		},
		LogErrors: true,
		CustomFormat: func(c *fiber.Ctx, duration time.Duration, statusCode int) string {
			return fmt.Sprintf("[%s] %s %d - %v - User-Agent: %s",
				c.Method(),
				c.Path(),
				statusCode,
				duration,
				c.Get("User-Agent"),
			)
		},
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "error",
				"database": "disconnected",
				"error":    err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"message":  "Service is healthy",
		})
	})

	// Database info endpoint
	app.Get("/db/info", func(c *fiber.Ctx) error {
		var version string
		err := db.QueryRow("SELECT version()").Scan(&version)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var currentDB string
		err = db.QueryRow("SELECT current_database()").Scan(&currentDB)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"version":          version,
			"current_database": currentDB,
			"status":           "connected",
		})
	})

	// Routes
	api := app.Group("/api/v1")
	server.NewRoutes(db, *cfg, api)

	port := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(port))
}
