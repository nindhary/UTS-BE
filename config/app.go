package config

import (
	"crud-app/database"
	"crud-app/middleware"
	"crud-app/route"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func SetupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func NewApp() *fiber.App {
	SetupLogger()
	database.ConnectDB()

	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)

	route.RegisterRoutes(app)
	return app
}
