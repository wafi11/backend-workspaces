package products

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func NewRoute(db *sql.DB, app fiber.Router) {
	repo := NewRepository(db)
	service := NewProductServices(repo)
	handler := NewHandler(service)

	// products routes
	products := app.Group("/products")
	products.Post("", handler.Create)
	products.Get("", handler.FindAll)
	products.Get("/:id", handler.FindById)
	products.Put("/:id", handler.Update)
	products.Delete("/:id", handler.Delete)

	container := app.Group("/container")
	container.Post("", handler.CreateContainer)
	container.Get("", handler.FindAllContainer)
	container.Get("/:id", handler.GetContainerByID)
	container.Put("/:id", handler.UpdateContainer)
	container.Delete("/:id", handler.DeleteContainer)

}
