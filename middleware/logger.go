package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	log.Printf("%s %s %d %v", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))
	return err
}
