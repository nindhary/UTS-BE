package service

import (
	models "crud-app/app/model"
	"crud-app/app/repository"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniService interface {
	GetAll() ([]models.Alumni, error)
	GetByID(id int) (models.Alumni, error)
	Create(input models.CreateAlumniRequest) (models.Alumni, error)
	Update(id int, input models.UpdateAlumniRequest) error
	Delete(id int) error
	// http handler
	GetAllHandler(c *fiber.Ctx) error
	GetByIDHandler(c *fiber.Ctx) error
	CreateHandler(c *fiber.Ctx) error
	UpdateHandler(c *fiber.Ctx) error
	DeleteHandler(c *fiber.Ctx) error
}

type alumniService struct {
	repo repository.AlumniRepository
}

func NewAlumniService(repo repository.AlumniRepository) AlumniService {
	return &alumniService{repo}
}

func (s *alumniService) GetAll() ([]models.Alumni, error) {
	return s.repo.GetAll()
}

func (s *alumniService) GetByID(id int) (models.Alumni, error) {
	return s.repo.GetByID(id)
}

func (s *alumniService) Create(input models.CreateAlumniRequest) (models.Alumni, error) {
	return s.repo.Create(input)
}

func (s *alumniService) Update(id int, input models.UpdateAlumniRequest) error {
	return s.repo.Update(id, input)
}

func (s *alumniService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *alumniService) GetAllHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "desc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	data, err := s.repo.GetAlumniRepo(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	total, _ := s.repo.CountAlumniRepo(search)

	response := models.AlumniResponse{
		Data: data,
		Meta: models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response)
}

func (s *alumniService) GetByIDHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Alumni tidak ditemukan",
				"data":    nil,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "OK",
		"data":    data,
	})
}

func (s *alumniService) CreateHandler(c *fiber.Ctx) error {
	var req models.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Body tidak valid",
			"data":    nil,
		})
	}
	data, err := s.repo.Create(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal tambah alumni",
			"data":    nil,
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Alumni ditambahkan",
		"data":    data,
	})
}

func (s *alumniService) UpdateHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req models.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Body tidak valid",
			"data":    nil,
		})
	}
	err := s.repo.Update(id, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal update alumni",
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni diupdate",
	})
}

func (s *alumniService) DeleteHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := s.repo.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal hapus alumni",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni dihapus",
	})
}
