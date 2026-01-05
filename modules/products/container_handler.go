package products

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wafi11/backend-workspaces/pkg/response"
)

func (h *Handlers) CreateContainer(c *fiber.Ctx) error {
	var req Container

	err := c.BodyParser(&req)

	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid body request")
	}

	cont, err := h.s.CreateContainer(c.Context(), req)

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to create container")
	}

	return response.Success(c, http.StatusCreated, "create container successfully", cont)

}

func (h *Handlers) GetContainerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return response.Error(c, http.StatusBadRequest, "id must be number")
	}

	data, err := h.s.GetContainerByID(c.Context(), idInt)

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to get container")
	}

	return response.Success(c, http.StatusOK, "get data container successfully", data)
}

func (h *Handlers) FindAllContainer(c *fiber.Ctx) error {

	data, err := h.s.FindAll(c.Context())

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to get container")
	}

	return response.Success(c, http.StatusOK, "get data container successfully", data)
}

func (h *Handlers) UpdateContainer(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return response.Error(c, http.StatusBadRequest, "id must be number")
	}
	var req Container

	err = c.BodyParser(&req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid body request")
	}
	err = h.s.UpdateContainer(c.Context(), idInt, req)

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to update container")
	}
	return response.Success(c, http.StatusOK, "update container successfully", id)
}

func (h *Handlers) DeleteContainer(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return response.Error(c, http.StatusBadRequest, "id must be number")
	}

	err = h.s.DeleteContainer(c.Context(), idInt)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to delete container")
	}

	return response.Success(c, http.StatusOK, "delete container successfully", id)
}
