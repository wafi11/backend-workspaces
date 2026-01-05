package auth

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func NewAuthRoute(db *sql.DB, app fiber.Router) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	auth := app.Group("/auth")
	auth.Post("", handler.RegisterUser)
}
