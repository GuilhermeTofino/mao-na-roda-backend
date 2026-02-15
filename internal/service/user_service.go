package service

import (
	"context"
	"errors"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/google/uuid"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateProfile(ctx context.Context, user *domain.User) error {
	// Verifica se user ID é válido (deve vir do auth middleware/token)
	if user.ID == uuid.Nil {
		return errors.New("ID de usuário inválido")
	}

	// Verifica duplicidade básica (opcional, DB tem constraint PK)
	existing, _ := s.repo.GetByID(ctx, user.ID)
	if existing != nil {
		return errors.New("perfil já existente")
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
