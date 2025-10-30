package middleware

import (
	"context"
	"crud-app/app/repository"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token diperlukan"})
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Format token salah"})
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Hanya admin yang boleh"})
		}
		return c.Next()
	}
}

func UserOrAdmin(pekerjaanRepo repository.PekerjaanRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		userID := fmt.Sprint(c.Locals("user_id")) // user ID dari token JWT
		targetID := c.Params("id")                // id pekerjaan dari URL (misal /pekerjaan/:id)

		fmt.Println("🔹 role:", role)
		fmt.Println("🔹 userID dari token:", userID)
		fmt.Println("🔹 target pekerjaan ID:", targetID)

		// Admin boleh lanjut tanpa dicek
		if role == "admin" {
			return c.Next()
		}

		// Cek apakah pekerjaan dengan ID itu milik user yang sedang login
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pekerjaanID, err := strconv.Atoi(targetID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
		}

		data, err := pekerjaanRepo.GetByID(pekerjaanID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Data pekerjaan tidak ditemukan"})
		}

		// Bandingkan user ID dari token dengan user ID di data pekerjaan
		if fmt.Sprint(data.UserID) != userID {
			return c.Status(403).JSON(fiber.Map{"error": "Kamu tidak punya akses ke data ini"})
		}

		return c.Next()
	}
}
