package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	StatusCode int         `json:"status_code"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		StatusCode: statusCode,
		Success:    true,
		Message:    message,
		Data:       data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		StatusCode: statusCode,
		Success:    false,
		Message:    message,
		Data:       nil,
	})
}
