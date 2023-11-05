package actions

import (
	"GrooveGuru/db"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

const (
	week = time.Hour * 24 * 7
)

var client, ctx = db.Instance()

func logError(fn, context string, err error) {
	log.Printf(
		"[%s] [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		fn, context, err.Error(),
	)
}

func badRequest(c *fiber.Ctx, msg string, code ...int) error {
	status := fiber.StatusBadRequest
	if len(code) > 0 {
		status = code[0]
	}
	return c.Status(status).JSON(fiber.Map{
		"error":   "bad request",
		"message": msg,
	})
}

func internalServerError(c *fiber.Ctx, msg string, code ...int) error {
	status := fiber.StatusInternalServerError
	if len(code) > 0 {
		status = code[0]
	}
	return c.Status(status).JSON(fiber.Map{
		"error":   "internal server error",
		"message": msg,
	})
}
