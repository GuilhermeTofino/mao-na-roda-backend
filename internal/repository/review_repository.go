package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reviewRepo struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) domain.ReviewRepository {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) Create(ctx context.Context, review *domain.Review) error {
	query := `
		INSERT INTO avaliacoes (id, profissional_id, cliente_id, nota, comentario, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	review.CreatedAt = time.Now()

	// Inicia transação para garantir consistência ao atualizar a média do profissional
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query,
		review.ID, review.ProfessionalID, review.ClientID,
		review.Rating, review.Comment, review.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar avaliação: %w", err)
	}

	// Atualiza nota média do profissional
	// Lógica simplificada: Incrementa contagem e recalcula média
	updateQuery := `
		UPDATE professionals 
		SET rating = ((rating * review_count) + $1) / (review_count + 1),
		    review_count = review_count + 1
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateQuery, float64(review.Rating), review.ProfessionalID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar nota do profissional: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *reviewRepo) GetByService(ctx context.Context, serviceID uuid.UUID) (*domain.Review, error) {
	return nil, nil // TODO
}

func (r *reviewRepo) ListByProfessional(ctx context.Context, proID uuid.UUID) ([]*domain.Review, error) {
	return nil, nil // TODO
}
