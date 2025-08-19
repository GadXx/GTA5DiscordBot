package service

import (
	"discordbot/internal/models"
	"discordbot/internal/repository"
	"time"
)

type Service struct {
	repo *repository.Repository
}

func NewInactiveService(r *repository.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) SaveRequest(userID, userName, reason string, createdAt, endAt time.Time) error {
	return s.repo.Save(userID, userName, reason, createdAt, endAt)
}

func (s *Service) ListVacations() ([]models.Inactive, error) {
	return s.repo.List()
}
