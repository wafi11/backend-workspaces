package products

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wafi11/backend-workspaces/pkg/response"
)

type Handlers struct {
	s Service
}

func NewHandler(s Service) Handlers {
	return Handlers{s: s}
}

func (h *Handlers) Create(c *fiber.Ctx) error {
	var req Product

	err := c.BodyParser(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	_, err = h.s.CreateProduct(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create product")
	}

	return response.Success(c, fiber.StatusCreated, "Product created successfully", req)
}

func (h *Handlers) FindAll(c *fiber.Ctx) error {
	products, err := h.s.GetAllProducts(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products")
	}

	return response.Success(c, fiber.StatusOK, "Products retrieved successfully", products)
}

func (h *Handlers) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return response.Error(c, http.StatusBadRequest, "id must be number")
	}

	var req Product
	err = c.BodyParser(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	err = h.s.UpdateProduct(c.Context(), idInt, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products")
	}

	return response.Success(c, fiber.StatusOK, "Update product successfully", nil)
}

func (h *Handlers) FindById(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return response.Error(c, http.StatusBadRequest, "id must be number")
	}
	products, err := h.s.GetProductByID(c.Context(), idInt)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products")
	}

	return response.Success(c, fiber.StatusOK, "Product retrieved successfully", products)
}

func (h *Handlers) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return response.Error(c, http.StatusBadRequest, "id must be number")
	}
	err = h.s.DeleteProduct(c.Context(), idInt)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products")
	}

	return response.Success(c, fiber.StatusOK, "Delete Product successfully", nil)
}
