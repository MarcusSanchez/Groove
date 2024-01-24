package util

import (
	"github.com/gofiber/fiber/v2"
)

func BadRequest(c *fiber.Ctx, msg string, code ...int) error {
	status := fiber.StatusBadRequest
	if len(code) > 0 {
		status = code[0]
	}
	return c.Status(status).JSON(fiber.Map{
		"error":   "bad request",
		"message": msg,
	})
}

func InternalServerError(c *fiber.Ctx, msg string, code ...int) error {
	status := fiber.StatusInternalServerError
	if len(code) > 0 {
		status = code[0]
	}
	return c.Status(status).JSON(fiber.Map{
		"error":   "internal server error",
		"message": msg,
	})
}

func Unauthorized(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   "unauthorized",
		"message": msg,
	})
}

func Forbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error":   "forbidden",
		"message": msg,
	})
}
