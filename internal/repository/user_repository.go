package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userAnswer struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userAnswer{db: db}
}

func (r *userAnswer) Create(ctx context.Context, user *domain.User) error {
	// A tabela é 'perfis', e assume-se que o ID vem do Supabase Auth (auth.users)
	query := `
		INSERT INTO perfis (id, nome_completo, tipo_usuario, avatar_url, telefone, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	user.CreatedAt = time.Now()

	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Type, user.AvatarURL, user.Phone, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar perfil: %w", err)
	}
	return nil
}

func (r *userAnswer) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, nome_completo, tipo_usuario, avatar_url, telefone, created_at FROM perfis WHERE id = $1`
	var user domain.User
	// Scan fields adaptados para as colunas do DB
	// Nota: telefone e avatar_url podem ser NULL no DB, então idealmente usar sql.NullString ou pgx types
	// Para simplificar, assumimos que puxe string vazia se null ou tratamos erro

	// Como pgx scan em string não aceita null direto sem wrapper, vamos usar *string temporário
	var avatar, phone *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Type, &avatar, &phone, &user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("perfil não encontrado: %w", err)
	}

	if avatar != nil {
		user.AvatarURL = *avatar
	}
	if phone != nil {
		user.Phone = *phone
	}

	return &user, nil
}
