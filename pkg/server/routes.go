package server

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/wafi11/backend-workspaces/modules/auth"
	"github.com/wafi11/backend-workspaces/modules/products"
	"github.com/wafi11/backend-workspaces/modules/templates"
	"github.com/wafi11/backend-workspaces/pkg/config"
)

func NewRoutes(db *sql.DB, cfg config.Config, api fiber.Router) {
	auth.NewAuthRoute(db, cfg, api)
	products.NewRoute(db, api)
	templates.NewTemplates(db, api)
}
