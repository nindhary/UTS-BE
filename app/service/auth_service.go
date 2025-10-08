package service

import (
	models "crud-app/app/model"
	"crud-app/middleware"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	LoginHandler(c *fiber.Ctx) error
	ProfileHandler(c *fiber.Ctx) error
}

type authService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) AuthService {
	return &authService{db: db}
}

// LoginHandler
func (a *authService) LoginHandler(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	var user models.User
	var passwordHash string

	// Query user berdasar username atau email
	err := a.db.QueryRow(`
        SELECT id, username, email, password_hash, role, created_at
        FROM users
        WHERE username=$1 OR email=$1
    `, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Username/password salah"})
	}

	// Debug log opsional
	fmt.Printf("DEBUG username: %q, password len=%d\n", req.Username, len(req.Password))

	// Cek password
	if !middleware.CheckPasswordHash(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Username/password salah"})
	}

	// Generate JWT token
	token, _ := middleware.GenerateToken(user)
	return c.JSON(models.LoginResponse{User: user, Token: token})
}

// ProfileHandler
func (a *authService) ProfileHandler(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)
	username, _ := c.Locals("username").(string)
	role, _ := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}
