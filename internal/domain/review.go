package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Review map para 'avaliacoes'
type Review struct {
	ID             uuid.UUID `json:"id"`
	ProfessionalID uuid.UUID `json:"profissional_id"`
	ClientID       uuid.UUID `json:"cliente_id"`
	Rating         int       `json:"nota"`
	Comment        string    `json:"comentario"`
	CreatedAt      time.Time `json:"created_at"`
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	ListByProfessional(ctx context.Context, proID uuid.UUID) ([]*Review, error)
}

type ReviewService interface {
	CreateReview(ctx context.Context, review *Review) error
}
