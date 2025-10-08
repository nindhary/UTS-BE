package main

import (
	// "crud-app/app/handler"

	"crud-app/config"
	"crud-app/database"
	"crud-app/route"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Koneksi database
	database.ConnectDB()
	defer database.DB.Close()

	// Inisialisasi Fiber dengan ErrorHandler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	route.RegisterRoutes(app)

	// Debug: cek daftar route yang sudah terdaftar
	for _, r := range app.GetRoutes() {
		log.Println(r.Method, r.Path)
	}

	// Start server
	log.Println("Server berjalan di http://localhost:3000")
	log.Fatal(app.Listen("0.0.0.0:3000"))
}
