package templates

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wafi11/backend-workspaces/pkg/response"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) Handler {
	return Handler{s: s}
}

func (h Handler) Create(c *fiber.Ctx) error {
	var req CreateTemplateRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid body request")
	}

	err := h.s.Create(c.Context(), req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, "Create Templates Successfully", nil)
}

func (h Handler) List(c *fiber.Ctx) error {
	limit := c.QueryInt("limit")
	cursor := c.Query("cursor")

	data, err := h.s.List(c.Context(), ListTemplatesRequest{
		Limit:  limit,
		Cursor: &cursor,
	})
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to retreived list templates")
	}

	return response.Success(c, http.StatusOK, "successfully to retrieved list templates", data)
}

func (h Handler) FindById(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "id must be number")
	}

	data, err := h.s.FindById(c.Context(), idInt)
	if err != nil {
		log.Printf("failed to get template : %s", err.Error())

		return response.Error(c, http.StatusInternalServerError, "failed to retreived get by id templates")
	}

	return response.Success(c, http.StatusOK, "successfully to retrieved template", data)
}
