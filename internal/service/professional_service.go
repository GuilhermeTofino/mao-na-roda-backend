package service

import (
	"context"
	"errors"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/google/uuid"
)

type professionalService struct {
	repo domain.ProfessionalRepository
}

func NewProfessionalService(repo domain.ProfessionalRepository) domain.ProfessionalService {
	return &professionalService{repo: repo}
}

func (s *professionalService) RegisterConfigs(ctx context.Context, pro *domain.Professional) error {
	// Verificar se já existe perfil para este ID
	existing, _ := s.repo.GetByID(ctx, pro.ID)
	if existing != nil {
		return errors.New("perfil profissional já existente")
	}

	// ID do profissional é SAME as User ID (perfil)
	// Verificar campos obrigatórios
	if pro.ID == uuid.Nil {
		return errors.New("ID de usuário inválido")
	}
	if pro.CategoryID == uuid.Nil {
		return errors.New("categoria inválida")
	}

	pro.Verified = false
	pro.Rating = 0

	return s.repo.Create(ctx, pro)
}

func (s *professionalService) GetProfile(ctx context.Context, id uuid.UUID) (*domain.Professional, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *professionalService) Search(ctx context.Context, categoryIDStr string, rating float64, lat, long, radius float64) ([]*domain.Professional, error) {
	var catID uuid.UUID
	if categoryIDStr != "" {
		parsed, err := uuid.Parse(categoryIDStr)
		if err == nil {
			catID = parsed
		}
	}

	filters := domain.SearchFilters{
		CategoryID: catID,
		MinRating:  rating,
		Latitude:   lat,
		Longitude:  long,
		RadiusKM:   radius,
	}
	return s.repo.Search(ctx, filters)
}

func (s *professionalService) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.repo.ListCategories(ctx)
}
