package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wafi11/backend-workspaces/pkg/response"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) RegisterUser(c *fiber.Ctx) error {
	var req RegisterUser

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid body request")
	}

	err := h.s.Register(c.Context(), req)
	if err != nil {
		statusCode := determineStatusCode(err)
		return response.Error(c, statusCode, err.Error())
	}

	return response.Success(c, http.StatusCreated, "user registered successfully", nil)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginUser

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid body request")
	}

	access_token, err := h.s.Login(c.Context(), req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	return response.Success(c, http.StatusOK, "Login Successfully", access_token)
}
