package service

import (
	"crud-app/app/repository"
	"crud-app/helper"
	"github.com/gofiber/fiber/v2"
)

type NganggurService struct {
	repo *repository.NganggurRepository
}

func NewNganggurService(repo *repository.NganggurRepository) *NganggurService {
	return &NganggurService{repo: repo}
}

func (s *NganggurService) GetAll(c *fiber.Ctx) error {
	alumni, err := s.repo.GetAll()
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	return helper.ResponseJSON(c, 200, "Success", true, fiber.Map{
		"total_alumni_tanpa_pekerjaan": len(alumni),
		"list_alumni":                  alumni,
	})
}
