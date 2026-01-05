package server

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/wafi11/backend-workspaces/modules/auth"
	"github.com/wafi11/backend-workspaces/modules/products"
)

func NewRoutes(db *sql.DB, api fiber.Router) {
	auth.NewAuthRoute(db, api)
	products.NewRoute(db, api)
}
