package templates

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func NewTemplates(db *sql.DB, app fiber.Router) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	api := app.Group("/templates")
	api.Post("", handler.Create)
	api.Get("", handler.List)
	api.Get("/:id", handler.FindById)
}
