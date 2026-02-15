package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ServiceStatus string

const (
	ServiceOpen       ServiceStatus = "Aberto"
	ServiceInProgress ServiceStatus = "Em andamento"
	ServiceCompleted  ServiceStatus = "Concluído"
	ServiceCancelled  ServiceStatus = "Cancelado"
)

// ServiceRequest representa um pedido de serviço
type ServiceRequest struct {
	ID             uuid.UUID     `json:"id"`
	ClientID       uuid.UUID     `json:"client_id"`
	ProfessionalID uuid.UUID     `json:"professional_id"`
	Description    string        `json:"description"`
	Status         ServiceStatus `json:"status"`
	ScheduledFor   time.Time     `json:"scheduled_for"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type ServiceRepository interface {
	Create(ctx context.Context, req *ServiceRequest) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status ServiceStatus) error
	GetByID(ctx context.Context, id uuid.UUID) (*ServiceRequest, error)
	ListByClient(ctx context.Context, clientID uuid.UUID) ([]*ServiceRequest, error)
	ListByProfessional(ctx context.Context, proID uuid.UUID) ([]*ServiceRequest, error)
}

type ServiceService interface {
	RequestService(ctx context.Context, clientID, proID uuid.UUID, description string, date time.Time) (*ServiceRequest, error)
	UpdateStatus(ctx context.Context, requestID uuid.UUID, status ServiceStatus, actorID uuid.UUID) error
}
