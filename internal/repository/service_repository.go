package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceRepo struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) domain.ServiceRepository {
	return &serviceRepo{db: db}
}

func (r *serviceRepo) Create(ctx context.Context, req *domain.ServiceRequest) error {
	query := `
		INSERT INTO requests (id, client_id, professional_id, description, status, scheduled_for, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		req.ID, req.ClientID, req.ProfessionalID, req.Description,
		req.Status, req.ScheduledFor, req.CreatedAt, req.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar solicitação: %w", err)
	}
	return nil
}

func (r *serviceRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ServiceStatus) error {
	query := `UPDATE requests SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status: %w", err)
	}
	return nil
}

func (r *serviceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServiceRequest, error) {
	query := `SELECT id, client_id, professional_id, description, status, scheduled_for, created_at, updated_at FROM requests WHERE id = $1`
	var req domain.ServiceRequest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&req.ID, &req.ClientID, &req.ProfessionalID, &req.Description,
		&req.Status, &req.ScheduledFor, &req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("solicitação não encontrada: %w", err)
	}
	return &req, nil
}

func (r *serviceRepo) ListByClient(ctx context.Context, clientID uuid.UUID) ([]*domain.ServiceRequest, error) {
	// Implementação similar ao GetByID, retornando slice
	return nil, nil // TODO
}

func (r *serviceRepo) ListByProfessional(ctx context.Context, proID uuid.UUID) ([]*domain.ServiceRequest, error) {
	// Implementação similar ao GetByID, retornando slice
	return nil, nil // TODO
}
