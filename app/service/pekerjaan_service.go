package service

import (
	models "crud-app/app/model"
	"crud-app/app/repository"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService interface {
	GetAll() ([]models.PekerjaanAlumni, error)
	GetByID(id int) (models.PekerjaanAlumni, error)
	Create(input models.CreatePekerjaanRequest) (models.PekerjaanAlumni, error)
	Update(id int, input models.UpdatePekerjaanRequest) (models.PekerjaanAlumni, error)
	Delete(id int) error
	GetByAlumni(alumniID int) ([]models.PekerjaanAlumni, error)
	// HTTP HANDLER
	GetAllHandler(c *fiber.Ctx) error
	GetByIDHandler(c *fiber.Ctx) error
	CreateHandler(c *fiber.Ctx) error
	UpdateHandler(c *fiber.Ctx) error
	DeleteHandler(c *fiber.Ctx) error
	GetByAlumniHandler(c *fiber.Ctx) error
	SoftDelete(c *fiber.Ctx) error
	GetTrash(c *fiber.Ctx) error
	Restore(c *fiber.Ctx) error
	HardDelete(c *fiber.Ctx) error
	// SoftDeleteBulk(c *fiber.Ctx) error
}

type pekerjaanService struct {
	repo repository.PekerjaanRepository
}

func NewPekerjaanService(repo repository.PekerjaanRepository) PekerjaanService {
	return &pekerjaanService{repo: repo}
}

func (s *pekerjaanService) GetAll() ([]models.PekerjaanAlumni, error) {
	return s.repo.GetAll()
}

func (s *pekerjaanService) GetByID(id int) (models.PekerjaanAlumni, error) {
	return s.repo.GetByID(id)
}

func (s *pekerjaanService) Create(input models.CreatePekerjaanRequest) (models.PekerjaanAlumni, error) {
	return s.repo.Create(input)
}

func (s *pekerjaanService) Update(id int, input models.UpdatePekerjaanRequest) (models.PekerjaanAlumni, error) {
	return s.repo.Update(id, input)
}

func (s *pekerjaanService) Delete(id int) error {
	return s.repo.Delete(id)
}
func (s *pekerjaanService) GetByAlumni(alumniID int) ([]models.PekerjaanAlumni, error) {
	return s.repo.GetByAlumni(alumniID)
}

func (s *pekerjaanService) GetAllHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	data, err := s.repo.GetPekerjaanRepo(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	total, _ := s.repo.CountPekerjaanRepo(search)

	response := fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"pages":  (total + limit - 1) / limit,
			"sortBy": sortBy,
			"order":  order,
			"search": search,
		},
	}

	return c.JSON(response)
}

func (s *pekerjaanService) GetByIDHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	data, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *pekerjaanService) CreateHandler(c *fiber.Ctx) error {
	var input models.CreatePekerjaanRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}
	data, err := s.repo.Create(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": data})
}

func (s *pekerjaanService) UpdateHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var input models.UpdatePekerjaanRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}
	data, err := s.repo.Update(id, input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *pekerjaanService) DeleteHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := s.repo.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil dihapus"})
}
func (s *pekerjaanService) GetByAlumniHandler(c *fiber.Ctx) error {
	alumniID, _ := strconv.Atoi(c.Params("alumni_id"))
	data, err := s.repo.GetByAlumni(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal ambil data",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "OK",
		"data":    data,
	})
}

func (s *pekerjaanService) SoftDelete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "id tidak valid"})
	}

	userIDInt, _ := strconv.Atoi(fmt.Sprint(c.Locals("user_id")))
	role, _ := c.Locals("role").(string)

	fmt.Println("JWT user_id:", userIDInt)
	fmt.Println("Role:", role)

	existingData, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "pekerjaan not found"})
	}
	fmt.Println("AlumniID pekerjaan:", existingData.AlumniID)

	if role != "admin" && existingData.UserID != userIDInt {
		return c.Status(403).JSON(fiber.Map{"message": "bukan pekerjaanmu"})
	}

	if err := s.repo.SoftDelete(id, models.UpdatePekerjaanRequest{}); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "pekerjaan soft deleted"})
}

func (s *pekerjaanService) GetTrash(c *fiber.Ctx) error {
	data, err := s.repo.GetTrash()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal ambil data",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "OK",
		"data":    data,
	})
}

func (s *pekerjaanService) Restore(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = s.repo.Restore(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "data berhasil di-restore"})
}

func (s *pekerjaanService) HardDelete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = s.repo.HardDelete(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "data berhasil dihapus permanen"})
}
