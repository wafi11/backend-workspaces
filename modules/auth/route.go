package auth

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/wafi11/backend-workspaces/pkg/config"
)

func NewAuthRoute(db *sql.DB, cfg config.Config, app fiber.Router) {

	repo := NewRepository(db, cfg)
	service := NewService(repo)
	handler := NewHandler(service)

	auth := app.Group("/auth")
	auth.Post("/register", handler.RegisterUser)
	auth.Post("/login", handler.Login)
}
