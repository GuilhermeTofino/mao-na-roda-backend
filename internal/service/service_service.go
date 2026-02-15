package service

import (
	"context"
	"errors"
	"time"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/google/uuid"
)

type serviceService struct {
	repo       domain.ServiceRepository
	reviewRepo domain.ReviewRepository
}

func NewServiceService(repo domain.ServiceRepository, reviewRepo domain.ReviewRepository) domain.ServiceService {
	return &serviceService{repo: repo, reviewRepo: reviewRepo}
}

func (s *serviceService) RequestService(ctx context.Context, clientID, proID uuid.UUID, description string, date time.Time) (*domain.ServiceRequest, error) {
	req := &domain.ServiceRequest{
		ID:             uuid.New(),
		ClientID:       clientID,
		ProfessionalID: proID,
		Description:    description,
		Status:         domain.ServiceOpen,
		ScheduledFor:   date,
	}

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *serviceService) UpdateStatus(ctx context.Context, requestID uuid.UUID, status domain.ServiceStatus, actorID uuid.UUID) error {
	// Verificar permissões (quem pode mudar status?)
	req, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	// Regras básicas de transição de estado
	if req.Status == domain.ServiceCompleted {
		return errors.New("serviço já concluído")
	}

	// Validar actorID (se é o cliente ou o profissional dono do serviço)
	if req.ClientID != actorID && req.ProfessionalID != actorID {
		return errors.New("não autorizado")
	}

	return s.repo.UpdateStatus(ctx, requestID, status)
}
