package route

import (
	"crud-app/app/repository"
	"crud-app/app/service"
	"crud-app/middleware"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// ===== DATABASE =====
	dsn := "postgres://postgres:12345678@localhost:5432/alumni?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// ===== BASE GROUP =====
	api := app.Group("/TM4")

	// ===== PUBLIC =====
	authSvc := service.NewAuthService(db)
	api.Post("/login", authSvc.LoginHandler)

	// ===== PROTECTED =====
	protected := api.Group("", middleware.AuthRequired())

	protected.Get("/profile", authSvc.ProfileHandler)
	// ===== ALUMNI =====
	alumniRepo := repository.NewAlumniRepository()
	alumniService := service.NewAlumniService(alumniRepo)

	// alumni nganggur
	nganggurRepo := &repository.NganggurRepository{DB: db}
	nganggurService := service.NewNganggurService(nganggurRepo)

	alumni := protected.Group("/alumni")
	alumni.Get("/", alumniService.GetAllHandler)
	alumni.Get("/:id", alumniService.GetByIDHandler)
	alumni.Post("/", middleware.AdminOnly(), alumniService.CreateHandler)
	alumni.Put("/:id", middleware.AdminOnly(), alumniService.UpdateHandler)
	alumni.Delete("/:id", middleware.AdminOnly(), alumniService.DeleteHandler)
	alumni.Get("/nganggur", middleware.AuthRequired(), nganggurService.GetAll)

	// ===== PEKERJAAN =====
	pekerjaanRepo := repository.NewPekerjaanRepository()
	pekerjaanService := service.NewPekerjaanService(pekerjaanRepo)
	pekerjaan := protected.Group("/pekerjaan")

	pekerjaan.Get("/trash", pekerjaanService.GetTrash)
	pekerjaan.Put("/restore/:id", pekerjaanService.Restore)
	pekerjaan.Delete("/harddelete/:id", pekerjaanService.HardDelete)

	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), pekerjaanService.GetByAlumniHandler)
	pekerjaan.Get("/:id", pekerjaanService.GetByIDHandler)
	pekerjaan.Get("/", pekerjaanService.GetAllHandler)
	pekerjaan.Post("/", middleware.AdminOnly(), pekerjaanService.CreateHandler)
	pekerjaan.Put("/:id", middleware.AdminOnly(), pekerjaanService.UpdateHandler)
	pekerjaan.Delete("/:id", middleware.AdminOnly(), pekerjaanService.DeleteHandler)
	pekerjaan.Delete("/soft/:id", pekerjaanService.SoftDelete)

	// pekerjaan.Delete("/soft/permanent", middleware.AdminOnly(), pekerjaanService.SoftDeleteBulk)
}
